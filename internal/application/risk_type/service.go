package risk_type

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RiskTypeService struct {
	repo RiskTypePort
}

type RiskTypeOutput struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RiskTypeId   int       `json:"risk_type_id"`
	Description  string    `json:"description"`
	RiskCategory string    `json:"riskCategory"`
	RiskTypeCode string    `json:"riskTypeCode"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (rrt *RiskTypeService) GetRiskTypes(ctx context.Context) ([]RiskTypeOutput, error) {
	results, err := rrt.repo.GetRiskTypes(ctx)
	if err != nil {
		return nil, err
	}
	riskTypes := make([]RiskTypeOutput, len(results))

	for i, r := range results {
		riskTypes[i] = RiskTypeOutput{
			ID:           r.ID,
			Name:         r.Name,
			RiskTypeCode: r.RiskTypeCode,
			Description:  r.Description,
			RiskCategory: r.RiskCategory,
			RiskTypeId:   r.RiskTypeId,
			CreatedAt:    r.CreatedAt,
		}
	}
	return riskTypes, nil
}
func (rrt *RiskTypeService) CreateRiskType(ctx context.Context) error {
	err := rrt.repo.CreateRiskType(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewRiskTypeService(repo RiskTypePort) RiskTypeService {
	return RiskTypeService{repo: repo}
}
