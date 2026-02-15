package domain

type BrownCard struct {
	RegistrationNumber string `json:"registrationNumber"`
	Success            bool   `json:"success"`
	Message            string `json:"message"`
	BrownCardNumber    string `json:"brownCardNumber"`
	URL                string `json:"url"`
}
