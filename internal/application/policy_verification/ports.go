package policy_verification

import (
	"context"

	"github.com/godsent-code/midtools/internal/domain"
)

type PolicyVerificationPort interface {
	GetPolicyVerification(ctx context.Context, cars []string) ([]domain.PolicyVerification, error)
}
