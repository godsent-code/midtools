package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/godsent-code/midtools/configs"
	"github.com/godsent-code/midtools/internal/domain"
	"golang.org/x/time/rate"
)

type USSDCheckRepository struct {
	config configs.Config
}

type nicUSSDCheckResponse struct {
	USERID  string `json:"USERID"`
	MSISDN  string `json:"MSISDN"`
	MSG     string `json:"MSG"`
	MSGTYPE bool   `json:"MSGTYPE"`
}

func (r *USSDCheckRepository) GetUSSDCheck(ctx context.Context, cars []string) ([]domain.USSDChecker, error) {
	results := make([]domain.USSDChecker, 0, len(cars))
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

				result := domain.USSDChecker{
					RegistrationNumber: car,
				}

				payload := map[string]interface{}{
					"registrationNumber": car,
					"USERID":             "1",
					"MSISDN":             "8",
					"MSGTYPE":            false,
					"USERDATA":           car,
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
					r.config.ApiEndPoint+"/public-api/vehicle-insurance-ussd-check",
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

				var nicResp nicUSSDCheckResponse
				if err := json.Unmarshal(bodyBytes, &nicResp); err != nil {
					result.Success = false
					result.Message = "Invalid response from NIC"
					r.appendResult(&mu, &results, result)
					continue
				}

				// Success case
				if nicResp.MSG != "" {
					result.Success = true
					result.Message = nicResp.MSG
				} else {
					result.Success = false
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

func (r *USSDCheckRepository) appendResult(mu *sync.Mutex, results *[]domain.USSDChecker, result domain.USSDChecker) {
	mu.Lock()
	*results = append(*results, result)
	mu.Unlock()
}

func NewUSSDCheckerRepository(config configs.Config) *USSDCheckRepository {
	return &USSDCheckRepository{config: config}
}
