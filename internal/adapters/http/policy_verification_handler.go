package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/godsent-code/midtools/internal/application/policy_verification"
	"github.com/godsent-code/midtools/pkg"
)

type PolicyVerificationHandler struct {
	service policy_verification.PolicyVerificationService
}

func (pvh *PolicyVerificationHandler) GetPolicyVerifications(w http.ResponseWriter, r *http.Request) {
	var request PolicyVerificationRequest

	if err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&request); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	br := policy_verification.PolicyVerificationInput{
		Cars: request.Cars,
	}

	if err := br.Validate(); err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	results, err := pvh.service.GetPolicyVerifications(r.Context(), br)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteResponse(w, http.StatusOK, results)
}

func NewPolicyVerificationHandler(service policy_verification.PolicyVerificationService) *PolicyVerificationHandler {
	return &PolicyVerificationHandler{service: service}
}
