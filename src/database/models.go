// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ApiKey struct {
	ID         pgtype.UUID        `json:"id"`
	UserID     pgtype.UUID        `json:"user_id"`
	KeyHash    string             `json:"key_hash"`
	Name       string             `json:"name"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	LastUsedAt pgtype.Timestamptz `json:"last_used_at"`
}

type Connection struct {
	ID              pgtype.UUID        `json:"id"`
	UserID          pgtype.UUID        `json:"user_id"`
	ProviderID      string             `json:"provider_id"`
	EncryptedApiKey string             `json:"encrypted_api_key"`
	Name            string             `json:"name"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	DeletedAt       pgtype.Timestamptz `json:"deleted_at"`
}

type Log struct {
	ID               pgtype.UUID        `json:"id"`
	UserID           pgtype.UUID        `json:"user_id"`
	ModelID          pgtype.UUID        `json:"model_id"`
	ConnectionID     pgtype.UUID        `json:"connection_id"`
	RequestPayload   []byte             `json:"request_payload"`
	ResponsePayload  []byte             `json:"response_payload"`
	PromptTokens     pgtype.Int8        `json:"prompt_tokens"`
	CompletionTokens pgtype.Int8        `json:"completion_tokens"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	Type             string             `json:"type"`
}

type Model struct {
	ID              pgtype.UUID        `json:"id"`
	UserID          pgtype.UUID        `json:"user_id"`
	ConnectionID    pgtype.UUID        `json:"connection_id"`
	ProxyModelID    string             `json:"proxy_model_id"`
	ProviderModelID string             `json:"provider_model_id"`
	Thinking        bool               `json:"thinking"`
	ToolsUsage      bool               `json:"tools_usage"`
	PriceInput      pgtype.Numeric     `json:"price_input"`
	PriceOutput     pgtype.Numeric     `json:"price_output"`
	DeletedAt       pgtype.Timestamptz `json:"deleted_at"`
	Type            string             `json:"type"`
}

type Provider struct {
	ID        pgtype.UUID        `json:"id"`
	UserID    pgtype.UUID        `json:"user_id"`
	Name      string             `json:"name"`
	BaseUrl   string             `json:"base_url"`
	Type      string             `json:"type"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type User struct {
	ID           pgtype.UUID        `json:"id"`
	Username     string             `json:"username"`
	PasswordHash string             `json:"password_hash"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}
