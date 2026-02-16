package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/godsent-code/midtools/configs"
	"github.com/godsent-code/midtools/internal/domain"
	"golang.org/x/time/rate"
)

type PolicyVerificationRepository struct {
	config configs.Config
}

type nicPolicyVerificationResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ProductName string `json:"productName"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
	} `json:"data"`
}

func (r *PolicyVerificationRepository) GetPolicyVerification(ctx context.Context, cars []string) ([]domain.PolicyVerification, error) {
	results := make([]domain.PolicyVerification, 0, len(cars))
	var mu sync.Mutex

	workerCount := 5
	jobs := make(chan string)

	limiter := rate.NewLimiter(rate.Every(300*time.Millisecond), 2)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for car := range jobs {

				if err := limiter.Wait(ctx); err != nil {
					return
				}

				result := domain.PolicyVerification{
					RegistrationNumber: car,
				}

				payload := map[string]interface{}{
					"data": map[string]string{
						"registrationNumber": strings.TrimSpace(car),
					},
				}

				jsonData, err := json.Marshal(payload)
				if err != nil {
					result.Success = false
					result.Message = "Failed to marshal request"
					r.appendResult(&mu, &results, result)
					continue
				}

				req, err := http.NewRequestWithContext(
					ctx,
					http.MethodPost,
					r.config.ApiEndPoint+"/public-api/policy-verification",
					bytes.NewBuffer(jsonData),
				)
				if err != nil {
					result.Success = false
					result.Message = "Failed to create request"
					r.appendResult(&mu, &results, result)
					continue
				}

				req.Header.Set("Authorization", "x-api-key "+r.config.ApiKey)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					result.Success = false
					result.Message = err.Error()
					r.appendResult(&mu, &results, result)
					continue
				}

				bodyBytes, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				var nicResp nicPolicyVerificationResponse
				if err := json.Unmarshal(bodyBytes, &nicResp); err != nil {
					result.Success = false
					result.Message = "Invalid response from NIC"
					r.appendResult(&mu, &results, result)
					continue
				}

				// Success case
				if nicResp.Success {
					result.Success = true
					result.ProductName = nicResp.Data.ProductName
					result.StartDate = nicResp.Data.StartDate
					result.EndDate = nicResp.Data.EndDate
					result.Message = "policy generated successfully."
				} else {
					result.Success = false
					result.Message = "Failed to generate Policy"
				}

				r.appendResult(&mu, &results, result)
			}
		}()
	}

	for _, car := range cars {
		jobs <- car
	}
	close(jobs)

	wg.Wait()
	return results, nil
}

func (r *PolicyVerificationRepository) appendResult(mu *sync.Mutex, results *[]domain.PolicyVerification, result domain.PolicyVerification) {
	mu.Lock()
	*results = append(*results, result)
	mu.Unlock()
}

func NewPolicyVerificationRepository(config configs.Config) *PolicyVerificationRepository {
	return &PolicyVerificationRepository{config: config}
}
