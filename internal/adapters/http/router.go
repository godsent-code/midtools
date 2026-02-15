package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/godsent-code/midtools/internal/application/brown_card_service"
	"github.com/godsent-code/midtools/internal/application/sticker"
)

func NewRouter(
	service brown_card_service.BrownCard,
	stickerService sticker.StickerService,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	BrownCard := NewBrownCardHandler(service)
	Sticker := NewStickerHandler(stickerService)

	r.Post("/browncard", BrownCard.GetBrownCard)
	r.Post("/sticker", Sticker.GetSticker)

	return r

}
