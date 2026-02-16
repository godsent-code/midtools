package product

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProductService struct {
	repo ProductPorts
}

type ProductOutput struct {
	ID          uuid.UUID
	Name        string
	ProductCode string
	ProductId   int
	Description string
	CreatedAt   time.Time
}

func (ps *ProductService) CreateProducts(ctx context.Context) error {
	err := ps.repo.CreateProduct(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) GetProducts(ctx context.Context) ([]ProductOutput, error) {
	results, err := ps.repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	products := make([]ProductOutput, len(results))
	for i, result := range results {
		products[i] = ProductOutput{
			ID:          result.ID,
			Name:        result.Name,
			ProductCode: result.ProductCode,
			ProductId:   result.ProductId,
			Description: result.Description,
			CreatedAt:   result.CreatedAt,
		}
	}
	return products, nil
}

func NewProductService(repo ProductPorts) ProductService {
	return ProductService{repo: repo}
}
