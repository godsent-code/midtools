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

type BrownCardRepository struct {
	config configs.Config
}
type nicResponse struct {
	Success bool `json:"success"`

	Data struct {
		StatusCode      string `json:"statusCode"`
		BrownCardNumber string `json:"brownCardNumber"`
		URL             string `json:"url"`
	} `json:"data"`

	Message        string `json:"message"`
	HttpStatusCode int    `json:"httpStatusCode"`
}

func (bcr *BrownCardRepository) GetBrownCard(ctx context.Context, cars []string) ([]domain.BrownCard, error) {

	results := make([]domain.BrownCard, 0, len(cars))
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

				result := domain.BrownCard{
					RegistrationNumber: car,
				}

				payload := map[string]interface{}{
					"data": map[string]string{
						"registrationNumber": car,
					},
				}

				jsonData, err := json.Marshal(payload)
				if err != nil {
					result.Success = false
					result.Message = "Failed to marshal request"
					appendResult(&mu, &results, result)
					continue
				}

				req, err := http.NewRequestWithContext(
					ctx,
					http.MethodPost,
					bcr.config.ApiEndPoint+"/public-api/generate-browncard",
					bytes.NewBuffer(jsonData),
				)
				if err != nil {
					result.Success = false
					result.Message = "Failed to create request"
					appendResult(&mu, &results, result)
					continue
				}

				req.Header.Set("Authorization", "x-api-key "+bcr.config.ApiKey)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					result.Success = false
					result.Message = err.Error()
					appendResult(&mu, &results, result)
					continue
				}

				bodyBytes, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				var nicResp nicResponse
				if err := json.Unmarshal(bodyBytes, &nicResp); err != nil {
					result.Success = false
					result.Message = "Invalid response from NIC"
					appendResult(&mu, &results, result)
					continue
				}

				// Success case
				if nicResp.Success {
					result.Success = true
					result.BrownCardNumber = nicResp.Data.BrownCardNumber
					result.URL = nicResp.Data.URL
					result.Message = "Brown card generated successfully"
				} else {
					result.Success = false
					result.Message = nicResp.Message
				}

				appendResult(&mu, &results, result)
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

func appendResult(mu *sync.Mutex, results *[]domain.BrownCard, result domain.BrownCard) {
	mu.Lock()
	*results = append(*results, result)
	mu.Unlock()
}

func NewBrownCardRepository(config configs.Config) *BrownCardRepository {
	return &BrownCardRepository{config: config}
}
