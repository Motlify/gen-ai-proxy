package api

import (
	"context"
	"fmt"

	"gen-ai-proxy/src/config"
	"gen-ai-proxy/src/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	db  *database.Queries
	cfg *config.Config
}

func NewService(db *database.Queries, cfg *config.Config) (*Service, error) {
	s := &Service{
		db:  db,
		cfg: cfg,
	}
	return s, nil
}

func (s *Service) GetProviderFromDB(ctx context.Context, providerID pgtype.UUID, userID pgtype.UUID) (Provider, error) {
	dbProvider, err := s.db.GetProvider(ctx, database.GetProviderParams{
		ID:     providerID,
		UserID: userID,
	})
	if err != nil {
		return Provider{}, err
	}
	return Provider{
		ID:      dbProvider.ID,
		Name:    dbProvider.Name,
		BaseURL: dbProvider.BaseUrl,
		Type:    dbProvider.Type,
	}, nil
}

func (s *Service) GetModelFromDB(ctx context.Context, modelID pgtype.UUID, userID pgtype.UUID) (Model, error) {
	dbModel, err := s.db.GetModel(ctx, database.GetModelParams{
		ID:     modelID,
		UserID: userID,
	})
	if err != nil {
		return Model{}, err
	}
	var priceInput, priceOutput float64
	if dbModel.PriceInput.Valid {
		f, err := dbModel.PriceInput.Float64Value()
		if err != nil {
			return Model{}, fmt.Errorf("cannot convert PriceInput to float64: %w", err)
		}
		priceInput = f.Float64
	}
	if dbModel.PriceOutput.Valid {
		f, err := dbModel.PriceOutput.Float64Value()
		if err != nil {
			return Model{}, fmt.Errorf("cannot convert PriceOutput to float64: %w", err)
		}
		priceOutput = f.Float64
	}
	return Model{
		ID:              dbModel.ID,
		ConnectionID:    dbModel.ConnectionID,
		ProxyModelID:    dbModel.ProxyModelID,
		ProviderModelID: dbModel.ProviderModelID,
		Thinking:        dbModel.Thinking,
		ToolsUsage:      dbModel.ToolsUsage,
		PriceInput:      priceInput,
		PriceOutput:     priceOutput,
	}, nil
}