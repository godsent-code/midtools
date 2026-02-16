package domain

import (
	"time"

	"github.com/google/uuid"
)

type RiskType struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RiskTypeId   int       `json:"risk_type_id"`
	Description  string    `json:"description"`
	RiskCategory string    `json:"riskCategory"`
	RiskTypeCode string    `json:"riskTypeCode"`
	CreatedAt    time.Time `json:"createdAt"`
}
