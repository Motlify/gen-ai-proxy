package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gen-ai-proxy/src/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// mustNumeric converts a float64 to pgtype.Numeric, handling potential errors.
func mustNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	str := strconv.FormatFloat(f, 'f', -1, 64)
	if err := n.Scan(str); err != nil {
		panic(fmt.Sprintf("failed to scan float64 into pgtype.Numeric: %v", err))
	}
	return n
}

// CreateModel godoc
// @Summary Create a new model
// @Schemes
// @Description Create a new model for a specific connection.
// @Tags Models
// @Accept json
// @Produce json
// @Param model body Model true "Model details"
// @Success 201 {object} Model
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/models [post]
func (s *Service) CreateModel(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var req struct {
		ConnectionID    string  `json:"connection_id"`
		ProviderModelID string  `json:"provider_model_id"`
		ProxyModelID    string  `json:"proxy_model_id"`
		PriceInput      float64 `json:"price_input"`
		PriceOutput     float64 `json:"price_output"`
		Thinking        bool    `json:"thinking"`
		ToolsUsage      bool    `json:"tools_usage"`
		Type            string  `json:"type"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	if req.ConnectionID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "connection_id is required"})
	}

	if req.Type == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "type is required"})
	}

	connectionID, err := uuid.Parse(req.ConnectionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Connection ID"})
	}

	_, err = s.db.GetConnection(c.Request().Context(), database.GetConnectionParams{
		ID:     pgtype.UUID{Bytes: connectionID, Valid: true},
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Connection with ID %s not found for this user", req.ConnectionID)})
	}

	modelPK := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	createdModel, err := s.db.CreateModel(c.Request().Context(), database.CreateModelParams{
		ID:              modelPK,
		UserID:          userID,
		ConnectionID:    pgtype.UUID{Bytes: connectionID, Valid: true},
		ProxyModelID:    req.ProxyModelID,
		ProviderModelID: req.ProviderModelID,
		Thinking:        req.Thinking,
		ToolsUsage:      req.ToolsUsage,
		PriceInput:      mustNumeric(req.PriceInput),
		PriceOutput:     mustNumeric(req.PriceOutput),
		Type:            req.Type,
	})
	if err != nil {
		fmt.Println("Error creating model in DB:", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create model"})
	}

	priceInputFloat, _ := createdModel.PriceInput.Float64Value()
	priceOutputFloat, _ := createdModel.PriceOutput.Float64Value()

	resp := Model{
		ID:              createdModel.ID,
		ConnectionID:    createdModel.ConnectionID,
		ProviderModelID: createdModel.ProviderModelID,
		ProxyModelID:    createdModel.ProxyModelID,
		Thinking:        createdModel.Thinking,
		ToolsUsage:      createdModel.ToolsUsage,
		PriceInput:      priceInputFloat.Float64,
		PriceOutput:     priceOutputFloat.Float64,
		Type:            createdModel.Type,
	}

	return c.JSON(http.StatusCreated, resp)
}

// UpdateModel godoc
// @Summary Update a model
// @Schemes
// @Description Update an existing model by ID.
// @Tags Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Param model body Model true "Model details"
// @Success 200 {object} Model
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/models/{id} [put]
func (s *Service) UpdateModel(c echo.Context) error {
	modelIDStr := c.Param("id")
	modelID, err := uuid.Parse(modelIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Model ID"})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var req struct {
		ProviderModelID string  `json:"provider_model_id"`
		ProxyModelID    string  `json:"proxy_model_id"`
		PriceInput      float64 `json:"price_input"`
		PriceOutput     float64 `json:"price_output"`
		Thinking        bool    `json:"thinking"`
		ToolsUsage      bool    `json:"tools_usage"`
		Type            string  `json:"type"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	updatedModel, err := s.db.UpdateModel(c.Request().Context(), database.UpdateModelParams{
		ID:              pgtype.UUID{Bytes: modelID, Valid: true},
		UserID:          userID,
		ProxyModelID:    req.ProxyModelID,
		ProviderModelID: req.ProviderModelID,
		Thinking:        req.Thinking,
		ToolsUsage:      req.ToolsUsage,
		PriceInput:      mustNumeric(req.PriceInput),
		PriceOutput:     mustNumeric(req.PriceOutput),
		Type:            req.Type,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update model"})
	}

	priceInputFloat, _ := updatedModel.PriceInput.Float64Value()
	priceOutputFloat, _ := updatedModel.PriceOutput.Float64Value()

	resp := Model{
		ID:              updatedModel.ID,
		ConnectionID:    updatedModel.ConnectionID,
		ProviderModelID: updatedModel.ProviderModelID,
		ProxyModelID:    updatedModel.ProxyModelID,
		Thinking:        updatedModel.Thinking,
		ToolsUsage:      updatedModel.ToolsUsage,
		PriceInput:      priceInputFloat.Float64,
		PriceOutput:     priceOutputFloat.Float64,
		Type:            updatedModel.Type,
	}

	return c.JSON(http.StatusOK, resp)
}

// ListModels godoc
// @Summary List all models
// @Schemes
// @Description List all available models.
// @Tags Models
// @Accept json
// @Produce json
// @Success 200 {array} Model
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/models [get]
func (s *Service) ListModels(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	dbModels, err := s.db.ListModels(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve models"})
	}

	respModels := make([]Model, len(dbModels))
	for i, m := range dbModels {
		priceInputFloat, _ := m.PriceInput.Float64Value()
		priceOutputFloat, _ := m.PriceOutput.Float64Value()

		respModels[i] = Model{
			ID:              m.ID,
			ConnectionID:    m.ConnectionID,
			ProviderModelID: m.ProviderModelID,
			ProxyModelID:    m.ProxyModelID,
			Thinking:        m.Thinking,
			ToolsUsage:      m.ToolsUsage,
			PriceInput:      priceInputFloat.Float64,
			PriceOutput:     priceOutputFloat.Float64,
			Type:            m.Type,
		}
	}
	return c.JSON(http.StatusOK, respModels)
}

// SoftDeleteModel godoc
// @Summary Soft delete a model
// @Schemes
// @Description Soft delete an existing model by ID.
// @Tags Models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/models/{id} [delete]
func (s *Service) SoftDeleteModel(c echo.Context) error {
	modelIDStr := c.Param("id")
	log.Printf("SoftDeleteModel: Received ID string: %s", modelIDStr)
	modelID, err := uuid.Parse(modelIDStr)
	if err != nil {
		log.Printf("SoftDeleteModel: Error parsing ID '%s': %v", modelIDStr, err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Model ID"})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	err = s.db.SoftDeleteModel(c.Request().Context(), database.SoftDeleteModelParams{
		ID:     pgtype.UUID{Bytes: modelID, Valid: true},
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to soft delete model"})
	}

	return c.NoContent(http.StatusNoContent)
}
