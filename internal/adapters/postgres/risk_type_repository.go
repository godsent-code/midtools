package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/godsent-code/midtools/configs"
	"github.com/godsent-code/midtools/internal/adapters/postgres/sqlc"
	"github.com/godsent-code/midtools/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type RiskTypeRepository struct {
	q      *pgxpool.Pool
	config configs.Config
}

type riskTypeNICResponse struct {
	Success bool `json:"success"`
	Data    struct {
		RiskTypes []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			RiskCategory string `json:"riskCategory"`
			RiskTypeCode string `json:"riskTypeCode"`
			Description  string `json:"description"`
		} `json:"riskTypes"`
	} `json:"data"`
}

func (rtr *RiskTypeRepository) GetRiskTypes(ctx context.Context) ([]*domain.RiskType, error) {
	q := sqlc.New(rtr.q)
	results, err := q.GetRiskType(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error getting risk types")
		return nil, err
	}
	riskTypes := make([]*domain.RiskType, len(results))
	for i, result := range results {
		riskTypes[i] = &domain.RiskType{
			Description:  result.Description.String,
			Name:         result.Name,
			ID:           result.ID,
			RiskCategory: result.RiskCategory,
			RiskTypeCode: result.RiskTypeCode,
			RiskTypeId:   int(result.RiskTypeID),
			CreatedAt:    result.CreatedAt.Time,
		}
	}
	return riskTypes, nil
}

func (rtr *RiskTypeRepository) CreateRiskType(ctx context.Context) error {
	tx, err := rtr.q.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Error begin transaction")
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	q := sqlc.New(tx)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		rtr.config.ApiEndPoint+"/public-api/risk-types",
		bytes.NewBuffer([]byte{}),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error create request")
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	req.Header.Set("Authorization", "x-api-key "+rtr.config.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error execute request")
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var nicResp riskTypeNICResponse
	if err := json.Unmarshal(bodyBytes, &nicResp); err != nil {
		log.Error().Err(err).Msg("Error unmarshal response")
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	if !nicResp.Success {
		log.Error().Msg("Error getting products from NIC")
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting products from NIC")
	}
	IDs := make([]int32, len(nicResp.Data.RiskTypes))
	name := make([]string, len(nicResp.Data.RiskTypes))
	riskCategory := make([]string, len(nicResp.Data.RiskTypes))
	riskTypeCode := make([]string, len(nicResp.Data.RiskTypes))
	description := make([]string, len(nicResp.Data.RiskTypes))
	for i := range nicResp.Data.RiskTypes {
		ids, err := strconv.Atoi(nicResp.Data.RiskTypes[i].ID)
		if err != nil {
			log.Error().Err(err).Msg("Error convert id to int")
			continue
		}
		IDs[i] = int32(ids)
		name[i] = nicResp.Data.RiskTypes[i].Name
		riskCategory[i] = nicResp.Data.RiskTypes[i].RiskCategory
		riskTypeCode[i] = nicResp.Data.RiskTypes[i].RiskTypeCode
		description[i] = nicResp.Data.RiskTypes[i].Description
	}

	err = q.CreateRiskType(ctx, sqlc.CreateRiskTypeParams{
		Name:         name,
		RiskTypeID:   IDs,
		RiskCategory: riskCategory,
		RiskTypeCode: riskTypeCode,
		Description:  description,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error create products")
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Error commit transaction")
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
	}
	return nil

}

func NewRiskTypeRepository(q *pgxpool.Pool, config configs.Config) *RiskTypeRepository {
	return &RiskTypeRepository{q: q, config: config}
}
