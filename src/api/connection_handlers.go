package api

import (
	"encoding/base64"
	"net/http"
	"time"

	"gen-ai-proxy/src/database"
	"gen-ai-proxy/src/encryption"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)


type ConnectionResponse struct {
	ID        pgtype.UUID `json:"id"`
	Provider  string    `json:"provider"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ListConnectionsResponse struct {
	Connections []ConnectionResponse `json:"connections"`
}

type CreateConnectionRequest struct {
	Name       string `json:"name" binding:"required"`
	APIKey     string `json:"api_key" binding:"required"`
	ProviderID string `json:"provider_id" binding:"required"`
}

type UpdateConnectionRequest struct {
	Name   string `json:"name" binding:"required"`
	APIKey string `json:"api_key" binding:"required"`
}
// ListConnections godoc
// @Summary List all connections
// @Schemes
// @Description List all connections for the authenticated user.
// @Tags Connections
// @Accept json
// @Produce json
// @Success 200 {object} ListConnectionsResponse
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/connections [get]
func (s *Service) ListConnections(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	dbConnections, err := s.db.ListConnections(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list connections"})
	}

	connections := make([]ConnectionResponse, len(dbConnections))
	for i, dbConnection := range dbConnections {
		connections[i] = ConnectionResponse{
			ID:        dbConnection.ID,
			Provider:  dbConnection.ProviderID,
			Name:      dbConnection.Name,
			CreatedAt: dbConnection.CreatedAt.Time,
		}
	}

	return c.JSON(http.StatusOK, ListConnectionsResponse{Connections: connections})
}



// CreateConnection godoc
// @Summary Create a new connection
// @Schemes
// @Description Create a new connection for the authenticated user.
// @Tags Connections
// @Accept json
// @Produce json
// @Param connection body CreateConnectionRequest true "Connection details"
// @Success 201
// @Failure 400 {object} api.ErrorResponse
// @Failure 401 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/connections [post]
func (s *Service) CreateConnection(c echo.Context) error {
	var req CreateConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	providerUUID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ProviderID format"})
	}

	if _, err := s.db.GetProvider(c.Request().Context(), database.GetProviderParams{ID: pgtype.UUID{Bytes: providerUUID, Valid: true}, UserID: userID}); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ProviderID"})
	}

	decodedEncryptionKey, err := base64.StdEncoding.DecodeString(s.cfg.EncryptionKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to decode encryption key"})
	}

	if len(decodedEncryptionKey) != 32 {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Encryption key must be 32 bytes long after base64 decoding"})
	}

	encryptedAPIKey, err := encryption.Encrypt(decodedEncryptionKey, []byte(req.APIKey))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to encrypt api key"})
	}

	params := database.CreateConnectionParams{
		UserID:          userID,
		ProviderID:      req.ProviderID,
		EncryptedApiKey: encryptedAPIKey,
		Name:            req.Name,
	}

	dbConnection, err := s.db.CreateConnection(c.Request().Context(), params)
	if err != nil {
		c.Logger().Errorf("Failed to create connection: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create connection"})
	}

	return c.JSON(http.StatusCreated, ConnectionResponse{
		ID:        dbConnection.ID,
		Provider:  dbConnection.ProviderID,
		Name:      dbConnection.Name,
		CreatedAt: dbConnection.CreatedAt.Time,
	})
}


// DeleteConnection godoc
// @Summary Delete a connection
// @Schemes
// @Description Delete a connection for the authenticated user.
// @Tags Connections
// @Accept json
// @Produce json
// @Param id path int true "Connection ID"
// @Success 204 "No Content"
// @Failure 400 {object} api.ErrorResponse
// @Failure 401 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Security BearerAuth
// @Router /api/connections/{id} [delete]
func (s *Service) DeleteConnection(c echo.Context) error {
	idStr := c.Param("id")

	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Connection ID format"})
	}

	connectionID := pgtype.UUID{Bytes: parsedID, Valid: true}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	err = s.db.SoftDeleteConnection(c.Request().Context(), database.SoftDeleteConnectionParams{
		ID: connectionID,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete connection"})
	}

	return c.NoContent(http.StatusNoContent)
}
