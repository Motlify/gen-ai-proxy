
package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gen-ai-proxy/src/database"
	"gen-ai-proxy/src/encryption"
	"gen-ai-proxy/src/llm"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)


type OpenAIUsage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int   `json:"total_tokens"`
}

type OpenAIResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []Choice      `json:"choices"`
	Usage   OpenAIUsage   `json:"usage"`
}

type Choice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ProxyOpenAIChat godoc
// @Summary Proxy a chat completion request to OpenAI
// @Schemes
// @Description Proxy a chat completion request to the OpenAI API.
// @Tags Proxy
// @Accept json
// @Produce json
// @Param request body ChatCompletionRequest true "OpenAI Chat Completion Request"
// @Success 200 {object} object
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /v1/chat/completions [post]
func (s *Service) ProxyOpenAIChat(c echo.Context) error {
	var err error
	var jsonBody []byte

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var req ChatCompletionRequest
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

	// Only allow OpenAI providers for this endpoint
	if llm.ProviderType(provider.Type) != llm.ProviderOpenAI {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "This endpoint only supports OpenAI providers"})
	}

	decodedEncryptionKey, err := base64.StdEncoding.DecodeString(s.cfg.EncryptionKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to decode encryption key"})
	}

	if len(decodedEncryptionKey) != 32 {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Encryption key must be 32 bytes long after base64 decoding"})
	}

	decryptedAPIKey, err := encryption.Decrypt(decodedEncryptionKey, connection.EncryptedApiKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to decrypt API key"})
	}

	apiKey := string(decryptedAPIKey)

	openAIReq := make(map[string]interface{})
	openAIReq["model"] = model.ProviderModelID
	openAIReq["stream"] = req.Stream

	openAIReq["messages"] = req.Messages

	if req.Tools != nil {
		toolsBytes, err := json.Marshal(req.Tools)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to marshal tools"})
		}
		var toolsData interface{}
		if err := json.Unmarshal(toolsBytes, &toolsData); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to unmarshal tools"})
		}
		openAIReq["tools"] = toolsData
	}

	// Marshal the request body to JSON
	jsonBody, err = json.Marshal(openAIReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to marshal request body"})
	}

	requestURL := provider.BaseUrl + "/chat/completions"
	log.Printf("Proxying OpenAI request to: %s", requestURL)
	proxyReq, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create proxy request"})
	}

	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
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
		// Log the conversation after successful streaming
		go func() {
			var finalOpenAIResp OpenAIResponse
			promptTokens := 0
			completionTokens := 0

			if err := json.Unmarshal(responseBody.Bytes(), &finalOpenAIResp); err == nil {
				promptTokens = int(finalOpenAIResp.Usage.PromptTokens)
				completionTokens = int(finalOpenAIResp.Usage.CompletionTokens)
			} else {
				log.Printf("Error unmarshaling final OpenAI streaming response for token counts: %v", err)
			}

			pt, err := safeIntToInt64(promptTokens)
			if err != nil {
				log.Printf("Prompt token conversion error: %v", err)
				return
			}

			ct, err := safeIntToInt64(completionTokens)
			if err != nil {
				log.Printf("Completion token conversion error: %v", err)
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

		log.Printf("OpenAI API Response Status: %d", resp.StatusCode)
		log.Printf("OpenAI API Response Body: %s", string(respBody))

		var data interface{}
		if err := json.Unmarshal(respBody, &data); err != nil {
			log.Printf("Failed to unmarshal OpenAI proxy response: %v, Raw Body: %s", err, string(respBody))
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to unmarshal proxy response"})
		}

		var openAIResp OpenAIResponse
		promptTokens := 0
		completionTokens := 0

		if err := json.Unmarshal(respBody, &openAIResp); err == nil {
			promptTokens = int(openAIResp.Usage.PromptTokens)
			completionTokens = int(openAIResp.Usage.CompletionTokens)
		} else {
			log.Printf("Error unmarshaling final OpenAI streaming response for token counts: %v", err)
		}

		pt, err := safeIntToInt64(promptTokens)
		if err != nil {
			log.Printf("Prompt token conversion error: %v", err)
			return err
		}

		ct, err := safeIntToInt64(completionTokens)
		if err != nil {
			log.Printf("Completion token conversion error: %v", err)
			return err
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
