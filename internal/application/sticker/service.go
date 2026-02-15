package sticker

import (
	"context"
	"errors"
	"strings"

	"github.com/godsent-code/midtools/pkg"
)

type StickerService struct {
	repo StickerPort
}
type StickerInput struct {
	Cars string
}

type StickerOutput struct {
	Status        bool   `json:"statusCode"`
	StickerLink   string `json:"stickerLink"`
	StickerNumber string `json:"stickerNumber"`
	Message       string `json:"message"`
	CarNumber     string `json:"carNumber"`
}

func (bci *StickerInput) Validate() error {
	if strings.TrimSpace(bci.Cars) == "" {
		return errors.New("cars is required")
	}
	return nil
}

func (ss *StickerService) GetSticker(ctx context.Context, input StickerInput) ([]StickerOutput, error) {
	parts := strings.FieldsFunc(input.Cars, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t'
	})

	if len(parts) == 0 {
		return nil, errors.New("cars is required")
	}

	stickers := make([]StickerOutput, 0)

	correctCars := make([]string, 0)

	for i, _ := range parts {
		exists, num := pkg.ValidateGhanaLicensePlate(parts[i])
		if !exists {
			stickers = append(stickers, StickerOutput{
				Status:        false,
				StickerNumber: num,
				CarNumber:     parts[i],
				StickerLink:   num,
				Message:       num,
			})
		} else {
			correctCars = append(correctCars, parts[i])
		}
	}

	results, err := ss.repo.GetStickers(ctx, correctCars)
	if err != nil {
		return nil, err
	}

	for i, _ := range results {
		stickers = append(stickers, StickerOutput{
			Status:        results[i].Success,
			StickerLink:   results[i].StickerLink,
			CarNumber:     results[i].RegistrationNumber,
			StickerNumber: results[i].StickerNumber,
			Message:       results[i].Message,
		})
	}

	return stickers, nil
}

func NewStickerService(r StickerPort) StickerService {
	return StickerService{repo: r}
}
