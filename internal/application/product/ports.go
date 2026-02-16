package product

import (
	"context"

	"github.com/godsent-code/midtools/internal/domain"
)

type ProductPorts interface {
	CreateProduct(ctx context.Context) error
	GetProducts(ctx context.Context) ([]*domain.Product, error)
}
