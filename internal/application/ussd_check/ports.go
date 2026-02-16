package ussd_check

import (
	"context"

	"github.com/godsent-code/midtools/internal/domain"
)

type USSDCheckPort interface {
	GetUSSDCheck(ctx context.Context, cars []string) ([]domain.USSDChecker, error)
}
