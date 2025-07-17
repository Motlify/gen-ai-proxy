package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gen-ai-proxy/src/database"
	"gen-ai-proxy/src/llm"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type OllamaChatRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	Stream   bool                    `json:"stream"`
	Think    bool                    `json:"think,omitempty"`
}

type OllamaResponse struct {
	Model           string    `json:"model"`
	CreatedAt       string    `json:"created_at"`
	Message         Message   `json:"message"`
	Done            bool      `json:"done"`
	TotalDuration   int64     `json:"total_duration"`
	LoadDuration    int64     `json:"load_duration"`
	PromptEvalCount int32     `json:"prompt_eval_count"`
	EvalCount       int32     `json:"eval_eval_count"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ProxyOllamaChat godoc
// @Summary Proxy a chat request to Ollama
// @Schemes
// @Description Proxy a chat request to the Ollama API.
// @Tags Proxy
// @Accept json
// @Produce json
// @Param request body OllamaChatRequest true "Ollama Chat Request"
// @Success 200 {object} object
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/chat [post]
func (s *Service) ProxyOllamaChat(c echo.Context) error {
	var err error
	var jsonBody []byte

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var req OllamaChatRequest
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	model, err := s.db.GetModelByProxyModelID(c.Request().Context(), database.GetModelByProxyModelIDParams{
		ProxyModelID:   req.Model,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Model not found"})
	}

	connection, err := s.db.GetConnection(c.Request().Context(), database.GetConnectionParams{
		ID:     model.ConnectionID,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get connection"})
	}

	providerUUID, err := uuid.Parse(connection.ProviderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to parse provider id"})
	}

	provider, err := s.db.GetProvider(c.Request().Context(), database.GetProviderParams{
		ID:     pgtype.UUID{Bytes: providerUUID, Valid: true},
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Provider not found for model"})
	}

	// Only allow Ollama providers for this endpoint
	if llm.ProviderType(provider.Type) != llm.ProviderOllama {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "This endpoint only supports Ollama providers"})
	}

	// Build Ollama request structure
	ollamaReq := make(map[string]interface{})
	ollamaReq["model"] = model.ProviderModelID
	ollamaReq["messages"] = req.Messages
	ollamaReq["stream"] = req.Stream

	// Add think parameter if specified
	if req.Think {
		ollamaReq["think"] = true
	}

	// Marshal the request body to JSON
	jsonBody, err = json.Marshal(ollamaReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to marshal request body"})
	}

	requestURL := provider.BaseUrl + "/api/chat"
	log.Printf("Proxying Ollama request to: %s", requestURL)
	proxyReq, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create proxy request"})
	}

	proxyReq.Header.Set("Content-Type", "application/json")
	// Note: Ollama typically doesn't require Authorization header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error sending proxy request to Ollama: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to send proxy request"})
	}
	defer resp.Body.Close()

	// Capture response body for logging
	var responseBody bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &responseBody)

	if req.Stream {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
		c.Response().WriteHeader(http.StatusOK)

		buf := make([]byte, 4096)
		for {
			n, err := teeReader.Read(buf)
			if n > 0 {
				if _, writeErr := c.Response().Write(buf[:n]); writeErr != nil {
					// Log the conversation even if there's a write error to the client
					go func() {
						_, logErr := s.db.CreateConversationLog(context.Background(), database.CreateConversationLogParams{
							UserID:          userID,
							ModelID:         model.ID,
							RequestPayload:  json.RawMessage(jsonBody),
							ResponsePayload: json.RawMessage(responseBody.Bytes()),
						})
						if logErr != nil {
							log.Printf("Error logging conversation (write error path): %v", logErr)
						}
					}()
					return writeErr
				}
				c.Response().Flush()
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				// Log the conversation even if there's a read error from the proxy
				go func() {
					_, logErr := s.db.CreateConversationLog(context.Background(), database.CreateConversationLogParams{
						UserID:          userID,
						ModelID:         model.ID,
						RequestPayload:  json.RawMessage(jsonBody),
						ResponsePayload: json.RawMessage(responseBody.Bytes()),
					})
					if logErr != nil {
						log.Printf("Error logging conversation (read error path): %v", logErr)
					}
				}()
				return err
			}
		}
		// After streaming is complete, parse the full response body for token counts
		var finalOllamaResp OllamaResponse
		var promptTokens int
		var completionTokens int

		if err := json.Unmarshal(responseBody.Bytes(), &finalOllamaResp); err == nil {
			promptTokens = int(finalOllamaResp.PromptEvalCount)
			completionTokens = int(finalOllamaResp.EvalCount)
		} else {
			log.Printf("Error unmarshaling final Ollama streaming response for token counts: %v", err)
		}

		// Log the conversation after successful streaming
		go func() {
			pt, err := safeIntToInt64(promptTokens)
			if err != nil {
				log.Printf("Error converting promptTokens: %v", err)
				return
			}
			ct, err := safeIntToInt64(completionTokens)
			if err != nil {
				log.Printf("Error converting completionTokens: %v", err)
				return
			}

			_, logErr := s.db.CreateConversationLog(context.Background(), database.CreateConversationLogParams{
				UserID:           userID,
				ModelID:          model.ID,
				RequestPayload:   json.RawMessage(jsonBody),
				ResponsePayload:  json.RawMessage(responseBody.Bytes()),
				PromptTokens:     pgtype.Int8{Int64: pt, Valid: true},
				CompletionTokens: pgtype.Int8{Int64: ct, Valid: true},
				ConnectionID:     connection.ID,
			})
			if logErr != nil {
				log.Printf("Error logging conversation (streaming success path): %v", logErr)
			}
		}()
		return nil
	} else {
		respBody, err := io.ReadAll(teeReader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read proxy response body"})
		}

		log.Printf("Raw Ollama non-streaming response: %s", string(respBody))
		var data interface{}
		if err := json.Unmarshal(respBody, &data); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to unmarshal proxy response"})
		}

		var ollamaResp OllamaResponse
		var promptTokens int
		var completionTokens int

		if err := json.Unmarshal(respBody, &ollamaResp); err == nil {
			promptTokens = int(ollamaResp.PromptEvalCount)
			completionTokens = int(ollamaResp.EvalCount)
		}

		pt, err := safeIntToInt64(promptTokens)
		if err != nil {
			log.Printf("Error converting promptTokens: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to process tokens"})
		}
		ct, err := safeIntToInt64(completionTokens)
		if err != nil {
			log.Printf("Error converting completionTokens: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to process tokens"})
		}

		_, logErr := s.db.CreateConversationLog(context.Background(), database.CreateConversationLogParams{
			UserID:           userID,
			ModelID:          model.ID,
			RequestPayload:   json.RawMessage(jsonBody),
			ResponsePayload:  json.RawMessage(respBody),
			PromptTokens:     pgtype.Int8{Int64: pt, Valid: true},
			CompletionTokens: pgtype.Int8{Int64: ct, Valid: true},
			ConnectionID:     connection.ID,
		})
		if logErr != nil {
			log.Printf("Error logging conversation (non-streaming success path): %v", logErr)
		}

		return c.JSON(resp.StatusCode, data)
	}
}
