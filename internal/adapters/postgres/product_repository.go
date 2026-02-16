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

type ProductRepository struct {
	q      *pgxpool.Pool
	config configs.Config
}
type productNICResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Products []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ProductCode string `json:"productCode"`
			Description string `json:"description"`
		} `json:"products"`
	} `json:"data"`
}

func (pr *ProductRepository) CreateProduct(ctx context.Context) error {
	tx, err := pr.q.BeginTx(ctx, pgx.TxOptions{})
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
		pr.config.ApiEndPoint+"/public-api/products",
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
	req.Header.Set("Authorization", "x-api-key "+pr.config.ApiKey)
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
	var nicResp productNICResponse
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
	productIDs := make([]int32, len(nicResp.Data.Products))
	ProductName := make([]string, len(nicResp.Data.Products))
	ProductCode := make([]string, len(nicResp.Data.Products))
	ProductDescription := make([]string, len(nicResp.Data.Products))
	for i := range nicResp.Data.Products {
		ids, err := strconv.Atoi(nicResp.Data.Products[i].ID)
		if err != nil {
			log.Error().Err(err).Msg("Error convert id to int")
			continue
		}
		productIDs[i] = int32(ids)
		ProductName[i] = nicResp.Data.Products[i].Name
		ProductCode[i] = nicResp.Data.Products[i].ProductCode
		ProductDescription[i] = nicResp.Data.Products[i].Description
	}

	err = q.CreateProducts(ctx, sqlc.CreateProductsParams{
		Name:        ProductName,
		ProductCode: ProductCode,
		ProductID:   productIDs,
		Description: ProductDescription,
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

func (pr *ProductRepository) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	q := sqlc.New(pr.q)
	results, err := q.GetProducts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error get products")
		return nil, err
	}
	products := make([]*domain.Product, len(results))
	for i, result := range results {
		products[i] = &domain.Product{
			ID:          result.ID,
			Name:        result.Name,
			ProductCode: result.ProductCode,
			ProductId:   int(result.ProductID),
			Description: result.Description.String,
			CreatedAt:   result.CreatedAt.Time,
		}
	}
	return products, nil
}

func NewProductRepository(pool *pgxpool.Pool, config configs.Config) *ProductRepository {
	return &ProductRepository{q: pool, config: config}
}
