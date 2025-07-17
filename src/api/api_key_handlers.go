package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"crypto/sha256"

	"gen-ai-proxy/src/database"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

type ListAPIKeysResponse struct {
	APIKeys []APIKeyResponse `json:"api_keys"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateAPIKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateAPIKeyResponse struct {
	APIKey string `json:"api_key"`
	Name   string `json:"name"`
}

type APIKeyResponse struct {
	ID         pgtype.UUID `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}


// ListAPIKeys godoc
// @Summary List all API keys
// @Schemes
// @Description List all API keys for the authenticated user.
// @Tags API Keys
// @Accept json
// @Produce json
// @Success 200 {object} ListAPIKeysResponse
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/api-keys [get]
func (s *Service) ListAPIKeys(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	dbAPIKeys, err := s.db.ListAPIKeys(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list api keys"})
	}

	apiKeys := make([]APIKeyResponse, len(dbAPIKeys))
	for i, dbAPIKey := range dbAPIKeys {
		// log.Printf("ListAPIKeys: Processing API Key ID: %v, Name: %s", dbAPIKey.ID, dbAPIKey.Name) // Debug log
		apiKeys[i] = APIKeyResponse{
			ID:         dbAPIKey.ID,
			Name:       dbAPIKey.Name,
			CreatedAt:  func() time.Time {
				if dbAPIKey.CreatedAt.Valid {
					return dbAPIKey.CreatedAt.Time
				}
				return time.Time{}
			}(),
			LastUsedAt: func() time.Time {
				if dbAPIKey.LastUsedAt.Valid {
					return dbAPIKey.LastUsedAt.Time
				}
				return time.Time{}
			}(),
		}
	}

	return c.JSON(http.StatusOK, ListAPIKeysResponse{APIKeys: apiKeys})
}

// CreateAPIKey godoc
// @Summary Create a new API key
// @Schemes
// @Description Create a new API key for the authenticated user.
// @Tags API Keys
// @Accept json
// @Produce json
// @Param key body CreateAPIKeyRequest true "API key details"
// @Success 201 {object} CreateAPIKeyResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/api-keys [post]
func (s *Service) CreateAPIKey(c echo.Context) error {
	var req CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate api key"})
	}
	apiKey := hex.EncodeToString(apiKeyBytes)

	hash := sha256.Sum256([]byte(apiKey))
	apiKeyHash := hex.EncodeToString(hash[:])

	params := database.CreateAPIKeyParams{
		UserID:  userID,
		KeyHash: apiKeyHash,
		Name:    req.Name,
	}

	dbAPIKey, err := s.db.CreateAPIKey(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create api key"})
	}

	return c.JSON(http.StatusCreated, CreateAPIKeyResponse{
		APIKey: apiKey,
		Name:   dbAPIKey.Name,
	})
}


// DeleteAPIKey godoc
// @Summary Delete an API key
// @Schemes
// @Description Delete an API key for the authenticated user.
// @Tags API Keys
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 204 "No Content"
// @Failure 400 {object} api.ErrorResponse
// @Failure 401 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/api-keys/{id} [delete]
func (s *Service) DeleteAPIKey(c echo.Context) error {
	idStr := c.Param("id")
	log.Printf("DeleteAPIKey: Received ID string: %s", idStr) // Debug log

	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("DeleteAPIKey: Error parsing ID '%s': %v", idStr, err) // Debug log
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid API Key ID format"})
	}

	apiKeyID := pgtype.UUID{Bytes: parsedID, Valid: true}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	err = s.db.DeleteAPIKey(c.Request().Context(), database.DeleteAPIKeyParams{
		ID: apiKeyID,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete api key"})
	}

	return c.NoContent(http.StatusNoContent)
}
