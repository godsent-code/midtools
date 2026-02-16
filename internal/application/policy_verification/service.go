package policy_verification

import (
	"context"
	"errors"
	"strings"

	"github.com/godsent-code/midtools/pkg"
)

type PolicyVerificationService struct {
	repo PolicyVerificationPort
}
type PolicyVerificationInput struct {
	Cars string
}

type PolicyVerificationOutput struct {
	Status      bool   `json:"statusCode"`
	ProductName string `json:"ProductName"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Message     string `json:"message"`
	CarNumber   string `json:"carNumber"`
}

func (bci *PolicyVerificationInput) Validate() error {
	if strings.TrimSpace(bci.Cars) == "" {
		return errors.New("cars is required")
	}
	return nil
}

func (pvs *PolicyVerificationService) GetPolicyVerifications(ctx context.Context, input PolicyVerificationInput) ([]PolicyVerificationOutput, error) {
	parts := strings.FieldsFunc(input.Cars, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t'
	})

	if len(parts) == 0 {
		return nil, errors.New("cars is required")
	}

	policies := make([]PolicyVerificationOutput, 0)

	correctCars := make([]string, 0)

	for i, _ := range parts {
		exists, num := pkg.ValidateGhanaLicensePlate(parts[i])
		if !exists {
			policies = append(policies, PolicyVerificationOutput{
				Status:    false,
				StartDate: num,
				CarNumber: parts[i],
				EndDate:   num,
				Message:   num,
			})
		} else {
			correctCars = append(correctCars, parts[i])
		}
	}

	results, err := pvs.repo.GetPolicyVerification(ctx, correctCars)
	if err != nil {
		return nil, err
	}

	for i, _ := range results {
		policies = append(policies, PolicyVerificationOutput{
			Status:      results[i].Success,
			StartDate:   results[i].StartDate,
			ProductName: results[i].ProductName,
			CarNumber:   results[i].RegistrationNumber,
			EndDate:     results[i].EndDate,
			Message:     results[i].Message,
		})
	}

	return policies, nil
}

func NewPolicyVerificationService(r PolicyVerificationPort) PolicyVerificationService {
	return PolicyVerificationService{repo: r}
}
