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

type OpenAIEmbedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type OpenAIEmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type OpenAIEmbeddingResponse struct {
	Object string               `json:"object"`
	Data   []OpenAIEmbedding    `json:"data"`
	Model  string               `json:"model"`
	Usage  OpenAIEmbeddingUsage `json:"usage"`
}

// ProxyOpenAIEmbedding godoc
// @Summary Proxy embedding request to OpenAI Compatible endpoint
// @Schemes
// @Description Proxy embedding request to the OpenAI API.
// @Tags Proxy
// @Accept json
// @Produce json
// @Param request body EmbeddingRequest true "OpenAI Embedding Request"
// @Success 200 {object} object
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /v1/embeddings [post]
func (s *Service) ProxyOpenAIEmbedding(c echo.Context) error {
	var err error
	var jsonBody []byte

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var req EmbeddingRequest
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	model, err := s.db.GetModelByProxyModelID(c.Request().Context(), database.GetModelByProxyModelIDParams{
		ProxyModelID: req.Model,
		UserID:       userID,
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

	if llm.ProviderType(provider.Type) != llm.ProviderOpenAI {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "This endpoint only supports OpenAI providers"})
	}

	if model.Type != "embedding" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "This endpoint only supports embedding models"})
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

	openAIReq := make(map[string]any)
	openAIReq["model"] = model.ProviderModelID
	openAIReq["input"] = req.Input
	if req.EncodingFormat != "" {
		openAIReq["encoding_format"] = req.EncodingFormat
	}

	jsonBody, err = json.Marshal(openAIReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to marshal request body"})
	}

	requestURL := provider.BaseUrl + "/embeddings"
	log.Printf("Proxying OpenAI embedding request to: %s", requestURL)
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read proxy response body"})
	}

	log.Printf("OpenAI API Response Status: %d", resp.StatusCode)
	log.Printf("OpenAI API Response Body: %s", string(respBody))

	var data any
	if err := json.Unmarshal(respBody, &data); err != nil {
		log.Printf("Failed to unmarshal OpenAI proxy response: %v, Raw Body: %s", err, string(respBody))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to unmarshal proxy response"})
	}

	var openAIResp OpenAIEmbeddingResponse
	promptTokens := 0

	if err := json.Unmarshal(respBody, &openAIResp); err == nil {
		promptTokens = openAIResp.Usage.PromptTokens
	} else {
		log.Printf("Error unmarshaling OpenAI embedding response for token counts: %v", err)
	}

	pt := int64(promptTokens)

	_, logErr := s.db.CreateLog(context.Background(), database.CreateLogParams{
		UserID:           userID,
		ModelID:          model.ID,
		RequestPayload:   json.RawMessage(jsonBody),
		ResponsePayload:  json.RawMessage(respBody),
		PromptTokens:     pgtype.Int8{Int64: pt, Valid: true},
		CompletionTokens: pgtype.Int8{Int64: 0, Valid: true}, // Embeddings does not generate completion tokens
		ConnectionID:     model.ConnectionID,
		Type:             "embedding",
	})
	if logErr != nil {
		log.Printf("Error logging embedding request: %v", logErr)
	}

	return c.JSON(resp.StatusCode, data)
}
