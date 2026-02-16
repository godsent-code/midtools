package risk_type

import (
	"context"

	"github.com/godsent-code/midtools/internal/domain"
)

type RiskTypePort interface {
	GetRiskTypes(ctx context.Context) ([]*domain.RiskType, error)
	CreateRiskType(ctx context.Context) error
}
