package api

import (
	"net/http"
	"log"

	"gen-ai-proxy/src/database"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

// CreateProvider godoc
// @Summary Create a new provider
// @Schemes
// @Description Create a new provider.
// @Tags Providers
// @Accept json
// @Produce json
// @Param provider body Provider true "Provider details"
// @Success 201 {object} Provider
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/providers [post]
func (s *Service) CreateProvider(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var reqMap map[string]interface{}
	if err := c.Bind(&reqMap); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	name, _ := reqMap["name"].(string)
	baseURL, _ := reqMap["base_url"].(string)
	providerType, _ := reqMap["type"].(string)

	log.Printf("CreateProvider: Received base_url: %s", baseURL)

	if providerType == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Provider type cannot be empty"})
	}

	providerID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	createdProvider, err := s.db.CreateProvider(c.Request().Context(), database.CreateProviderParams{
		ID:      providerID,
		UserID:  userID,
		Name:    name,
		BaseUrl: baseURL,
		Type:    providerType,
	})
	if err != nil {
		log.Printf("CreateProvider: Failed to create provider in DB: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create provider"})
	}

	resp := Provider{
		ID:      createdProvider.ID,
		Name:    createdProvider.Name,
		BaseURL: createdProvider.BaseUrl,
		Type:    createdProvider.Type,
	}
	log.Printf("CreateProvider: Created provider with BaseURL: %s", resp.BaseURL)

	return c.JSON(http.StatusCreated, resp)
}



// ListProviders godoc
// @Summary List all providers
// @Schemes
// @Description List all available providers.
// @Tags Providers
// @Accept json
// @Produce json
// @Param provider body Provider true "Provider details"
// @Success 200 {array} Provider
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/providers [get]
func (s *Service) ListProviders(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	dbProviders, err := s.db.ListProviders(c.Request().Context(), userID)
	if err != nil {
		log.Printf("ListProviders: Failed to retrieve providers from DB: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve providers"})
	}

	respProviders := make([]Provider, len(dbProviders))
	for i, p := range dbProviders {
		provider, err := s.GetProviderFromDB(c.Request().Context(), p.ID, userID)
		if err != nil {
			log.Printf("ListProviders: Failed to retrieve provider from DB: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve provider from DB"})
		}
		respProviders[i] = provider
		log.Printf("ListProviders: Retrieved provider %s with BaseURL: %s", provider.Name, provider.BaseURL)
	}
	return c.JSON(http.StatusOK, respProviders)
}

// DeleteProvider godoc
// @Summary Delete a provider
// @Schemes
// @Description Delete a provider by ID.
// @Tags Providers
// @Accept json
// @Produce json
// @Param id path string true "Provider ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/providers/{id} [delete]
func (s *Service) DeleteProvider(c echo.Context) error {
	providerIDStr := c.Param("id")
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Provider ID"})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	err = s.db.SoftDeleteProvider(c.Request().Context(), database.SoftDeleteProviderParams{
		ID:     pgtype.UUID{Bytes: providerID, Valid: true},
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to soft delete provider"})
	}

	// Soft delete all connections associated with this provider
	connections, err := s.db.ListConnectionsByProviderID(c.Request().Context(), database.ListConnectionsByProviderIDParams{
		ProviderID: providerIDStr,
		UserID:     userID,
	})
	if err != nil {
		log.Printf("Error listing connections for provider %s: %v", providerIDStr, err)
		// Continue with provider deletion even if connections cannot be listed
	}

	for _, conn := range connections {
		err := s.db.SoftDeleteConnection(c.Request().Context(), database.SoftDeleteConnectionParams{
			ID:     conn.ID,
			UserID: userID,
		})
		if err != nil {
			log.Printf("Error soft deleting connection %s for provider %s: %v", conn.ID.String(), providerIDStr, err)
			// Continue with other connections even if one fails
		}
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete provider"})
	}

	return c.NoContent(http.StatusNoContent)
}