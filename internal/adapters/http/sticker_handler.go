package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/godsent-code/midtools/internal/application/sticker"
	"github.com/godsent-code/midtools/pkg"
)

type StickerHandler struct {
	service sticker.StickerService
}

func (ach *StickerHandler) GetSticker(w http.ResponseWriter, r *http.Request) {
	var request BrownCardRequest

	if err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&request); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	br := sticker.StickerInput{
		Cars: request.Cars,
	}

	if err := br.Validate(); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	results, err := ach.service.GetSticker(r.Context(), br)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, results)
}

func NewStickerHandler(service sticker.StickerService) *StickerHandler {
	return &StickerHandler{service: service}
}
