package brown_card_service

import (
	"context"
	"errors"
	"strings"

	"github.com/godsent-code/midtools/pkg"
)

type BrownCard struct {
	repo BrownCardService
}

type BrownCardInput struct {
	Cars string
}
type BrownCardOutput struct {
	Status          bool   `json:"statusCode"`
	BrownCardNumber string `json:"brownCardNumber"`
	Url             string `json:"url"`
	Message         string `json:"message"`
	CarNumber       string `json:"carNumber"`
}

func (bci *BrownCardInput) Validate() error {
	if strings.TrimSpace(bci.Cars) == "" {
		return errors.New("cars is required")
	}
	return nil
}

func (bc *BrownCard) GetBrownCard(ctx context.Context, input BrownCardInput) ([]BrownCardOutput, error) {
	parts := strings.FieldsFunc(input.Cars, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t'
	})

	if len(parts) == 0 {
		return nil, errors.New("cars is required")
	}

	brownCards := make([]BrownCardOutput, 0)

	correctCars := make([]string, 0)

	for i, _ := range parts {
		exists, num := pkg.ValidateGhanaLicensePlate(parts[i])
		if !exists {
			brownCards = append(brownCards, BrownCardOutput{
				Status:          false,
				BrownCardNumber: num,
				CarNumber:       parts[i],
				Url:             num,
				Message:         num,
			})
		} else {
			correctCars = append(correctCars, parts[i])
		}
	}

	results, err := bc.repo.GetBrownCard(ctx, correctCars)
	if err != nil {
		return nil, err
	}

	for i, _ := range results {
		brownCards = append(brownCards, BrownCardOutput{
			Status:          results[i].Success,
			Url:             results[i].URL,
			CarNumber:       results[i].RegistrationNumber,
			BrownCardNumber: results[i].BrownCardNumber,
			Message:         results[i].Message,
		})
	}

	return brownCards, nil
}

func NewBrownCard(repo BrownCardService) BrownCard {
	return BrownCard{repo: repo}
}
