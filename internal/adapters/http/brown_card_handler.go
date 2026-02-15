package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/godsent-code/midtools/internal/application/brown_card_service"
	"github.com/godsent-code/midtools/pkg"
)

type BrownCardHandler struct {
	service brown_card_service.BrownCard
}

func (ach *BrownCardHandler) GetBrownCard(w http.ResponseWriter, r *http.Request) {
	var request BrownCardRequest

	if err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&request); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	br := brown_card_service.BrownCardInput{
		Cars: request.Cars,
	}

	if err := br.Validate(); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	results, err := ach.service.GetBrownCard(r.Context(), br)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, results)
}

func NewBrownCardHandler(service brown_card_service.BrownCard) *BrownCardHandler {
	return &BrownCardHandler{service: service}
}
