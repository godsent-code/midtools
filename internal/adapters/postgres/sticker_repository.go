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

type StickerRepository struct {
	config configs.Config
}

type nicStickerResponse struct {
	Success bool `json:"success"`

	Data struct {
		StickerLink   string `json:"stickerLink"`
		StickerNumber string `json:"stickerNumber"`
	} `json:"data"`

	Message string `json:"message"`
}

func (r *StickerRepository) GetStickers(ctx context.Context, cars []string) ([]domain.Sticker, error) {
	results := make([]domain.Sticker, 0, len(cars))
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

				result := domain.Sticker{
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
					r.appendResult(&mu, &results, result)
					continue
				}

				req, err := http.NewRequestWithContext(
					ctx,
					http.MethodPost,
					r.config.ApiEndPoint+"/public-api/generate-sticker",
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

				var nicResp nicStickerResponse
				if err := json.Unmarshal(bodyBytes, &nicResp); err != nil {
					result.Success = false
					result.Message = "Invalid response from NIC"
					r.appendResult(&mu, &results, result)
					continue
				}

				// Success case
				if nicResp.Success {
					result.Success = true
					result.StickerLink = nicResp.Data.StickerLink
					result.StickerNumber = nicResp.Data.StickerNumber
					result.Message = "Sticker generated and assigned to policy successfully."
				} else {
					result.Success = false
					result.Message = nicResp.Message
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

func (r *StickerRepository) appendResult(mu *sync.Mutex, results *[]domain.Sticker, result domain.Sticker) {
	mu.Lock()
	*results = append(*results, result)
	mu.Unlock()
}

func NewStickerRepository(config configs.Config) *StickerRepository {
	return &StickerRepository{config: config}
}
