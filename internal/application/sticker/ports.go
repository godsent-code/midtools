package sticker

import (
	"context"

	"github.com/godsent-code/midtools/internal/domain"
)

type StickerPort interface {
	GetStickers(ctx context.Context, cars []string) ([]domain.Sticker, error)
}
