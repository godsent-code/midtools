package ussd_check

import (
	"context"
	"errors"
	"strings"

	"github.com/godsent-code/midtools/pkg"
)

type USSDCheckService struct {
	repo USSDCheckPort
}
type USSDCheckInput struct {
	Cars string
}

type USSDCheckOutput struct {
	Status    bool   `json:"statusCode"`
	Message   string `json:"message"`
	CarNumber string `json:"carNumber"`
}

func (bci *USSDCheckInput) Validate() error {
	if strings.TrimSpace(bci.Cars) == "" {
		return errors.New("cars is required")
	}
	return nil
}

func (ss *USSDCheckService) GetUSSDCheck(ctx context.Context, input USSDCheckInput) ([]USSDCheckOutput, error) {
	parts := strings.FieldsFunc(input.Cars, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t'
	})

	if len(parts) == 0 {
		return nil, errors.New("cars is required")
	}

	ussds := make([]USSDCheckOutput, 0)

	correctCars := make([]string, 0)

	for i, _ := range parts {
		exists, num := pkg.ValidateGhanaLicensePlate(parts[i])
		if !exists {
			ussds = append(ussds, USSDCheckOutput{
				Status:    false,
				CarNumber: parts[i],
				Message:   num,
			})
		} else {
			correctCars = append(correctCars, parts[i])
		}
	}

	results, err := ss.repo.GetUSSDCheck(ctx, correctCars)
	if err != nil {
		return nil, err
	}

	for i, _ := range results {
		ussds = append(ussds, USSDCheckOutput{
			Status:    results[i].Success,
			CarNumber: results[i].RegistrationNumber,
			Message:   results[i].Message,
		})
	}

	return ussds, nil
}

func NewUSSDCheckService(r USSDCheckPort) USSDCheckService {
	return USSDCheckService{repo: r}
}
