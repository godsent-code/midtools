package domain

type USSDChecker struct {
	RegistrationNumber string `json:"registrationNumber"`
	Message            string `json:"message"`
	Status             string `json:"status"`
	Success            bool   `json:"success"`
}
