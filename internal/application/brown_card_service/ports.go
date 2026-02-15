package brown_card_service

import (
	"context"

	"github.com/godsent-code/midtools/internal/domain"
)

type BrownCardService interface {
	GetBrownCard(ctx context.Context, cars []string) ([]domain.BrownCard, error)
}
