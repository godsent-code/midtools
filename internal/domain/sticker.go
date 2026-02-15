package domain

type Sticker struct {
	StickerLink        string `json:"stickerLink"`
	StickerNumber      string `json:"stickerNumber"`
	Message            string `json:"message"`
	RegistrationNumber string `json:"registrationNumber"`
	Status             string `json:"status"`
	Success            bool   `json:"success"`
}
