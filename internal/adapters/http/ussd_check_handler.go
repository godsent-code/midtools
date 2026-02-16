package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/godsent-code/midtools/internal/application/ussd_check"
	"github.com/godsent-code/midtools/pkg"
)

type USSDCheckHandler struct {
	service ussd_check.USSDCheckService
}

func (usd *USSDCheckHandler) GetUSSDCheck(w http.ResponseWriter, r *http.Request) {
	var request USSDCheckRequest

	if err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&request); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	br := ussd_check.USSDCheckInput{
		Cars: request.Cars,
	}

	if err := br.Validate(); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	results, err := usd.service.GetUSSDCheck(r.Context(), br)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, results)
}

func NewUSSDCheckHandler(service ussd_check.USSDCheckService) *USSDCheckHandler {
	return &USSDCheckHandler{service: service}
}
