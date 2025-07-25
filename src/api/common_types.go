package api

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func safeIntToInt64(i int) (int64, error) {
	return int64(i), nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// OpenAI Compatible LLM Request
type ChatCompletionRequest struct {
	Model        string                  `json:"model"`
	ConnectionID string                  `json:"connection_id"`
	Messages     []ChatCompletionMessage `json:"messages"`
	Stream       bool                    `json:"stream,omitempty"`
	Tools        any                     `json:"tools,omitempty"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI Compatible Embedding Request

type EmbeddingRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"`
}

// API
type Provider struct {
	ID      pgtype.UUID `json:"id"`
	Name    string      `json:"name"`
	BaseURL string      `json:"base_url"`
	Type    string      `json:"type"`
}

type Model struct {
	ID              pgtype.UUID `json:"id"`
	ConnectionID    pgtype.UUID `json:"connection_id"`
	ProviderModelID string      `json:"provider_model_id"`
	ProxyModelID    string      `json:"proxy_model_id"`
	Thinking        bool        `json:"thinking"`
	ToolsUsage      bool        `json:"tools_usage"`
	PriceInput      float64     `json:"price_input"`
	PriceOutput     float64     `json:"price_output"`
	Type            string      `json:"type"`
}
