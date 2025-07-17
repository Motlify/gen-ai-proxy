package api

import (
	"net/http"
	"time"
	"gen-ai-proxy/src/database"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgtype"
)


type RawJSON string // @name RawJSON

type ConversationLogResponse struct {
	ID               pgtype.UUID     `json:"id"`
	ModelID          pgtype.UUID     `json:"model_id"`
	ConnectionID     pgtype.UUID     `json:"connection_id"`
	RequestPayload   RawJSON         `json:"request_payload"`
	ResponsePayload  RawJSON         `json:"response_payload"`
	CreatedAt        time.Time       `json:"created_at"`
	PromptTokens     int64           `json:"prompt_tokens"`
	CompletionTokens int64           `json:"completion_tokens"`
}

type ListConversationLogsRequest struct {
	Page        int64       `query:"page"`
	Limit       int64       `query:"limit"`
	ModelID     pgtype.UUID `query:"model_id"`
	ProviderID  pgtype.UUID `query:"provider_id"`
	ConnectionID pgtype.UUID `query:"connection_id"`
}

type ListConversationLogsResponse struct {
	Logs  []ConversationLogResponse `json:"logs"`
	Total int64                     `json:"total"`
}


// ListConversationLogs godoc
// @Summary List conversation logs
// @Schemes
// @Description List all conversation logs with pagination and filtering.
// @Tags Conversation Logs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param model_id query string false "Filter by model ID"
// @Param connection_id query int false "Filter by connection ID"
// @Success 200 {object} ListConversationLogsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/conversation_logs [get]
func (s *Service) ListConversationLogs(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
	}

	var req ListConversationLogsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	// Default pagination values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	params := database.ListConversationLogsParams{
		Limit:  pgtype.Int8{Int64: req.Limit, Valid: true},
		Offset: pgtype.Int8{Int64: offset, Valid: true},
		UserID: userID,
	}

	if req.ModelID.Valid {
		params.ModelID = req.ModelID
	}

	logs, err := s.db.ListConversationLogs(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve conversation logs"})
	}

	filteredLogs := []database.ListConversationLogsRow{}
	for _, log := range logs {
		// Filter by ProviderID
		if req.ProviderID.Valid {
			model, err := s.GetModelFromDB(c.Request().Context(), log.ModelID, userID)
			if err != nil {
				continue
			}
			connection, err := s.db.GetConnection(c.Request().Context(), database.GetConnectionParams{
				ID:     model.ConnectionID,
				UserID: userID,
			})
			if err != nil || (connection.ProviderID != req.ProviderID.String()) {
				continue
			}
		}

		// Filter by ConnectionID
		if req.ConnectionID.Valid {
			if !log.ConnectionID.Valid || (log.ConnectionID.Valid && req.ConnectionID.Valid && log.ConnectionID != req.ConnectionID) {
				continue
			}
		}
		filteredLogs = append(filteredLogs, log)
	}

	respLogs := make([]ConversationLogResponse, len(filteredLogs))
	for i, log := range filteredLogs {
		respLogs[i] = ConversationLogResponse{
			ID:             log.ID,
			ModelID:        log.ModelID,
			ConnectionID:   log.ConnectionID,
			RequestPayload: RawJSON(log.RequestPayload),
			ResponsePayload: RawJSON(log.ResponsePayload),
			CreatedAt:      log.CreatedAt.Time,
			PromptTokens:     func() int64 {
				if log.PromptTokens.Valid {
					return log.PromptTokens.Int64
				}
				return 0
			}(),
			CompletionTokens: func() int64 {
				if log.CompletionTokens.Valid {
					return log.CompletionTokens.Int64
				}
				return 0
			}(),
		}
	}

	return c.JSON(http.StatusOK, ListConversationLogsResponse{
		Logs:  respLogs,
		Total: int64(len(filteredLogs)),
	})
}

