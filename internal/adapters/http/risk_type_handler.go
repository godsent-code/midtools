package http

import (
	"net/http"

	"github.com/godsent-code/midtools/internal/application/risk_type"
	"github.com/godsent-code/midtools/pkg"
)

type RiskTypeHandler struct {
	service risk_type.RiskTypeService
}

func (rth *RiskTypeHandler) GetRiskTypes(w http.ResponseWriter, r *http.Request) {
	risks, err := rth.service.GetRiskTypes(r.Context())
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, risks)
}

func (rth *RiskTypeHandler) CreateRiskType(w http.ResponseWriter, r *http.Request) {
	err := rth.service.CreateRiskType(r.Context())
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, "Risk Type created")
}

func NewRiskTypeHandler(service risk_type.RiskTypeService) *RiskTypeHandler {
	return &RiskTypeHandler{service: service}
}
