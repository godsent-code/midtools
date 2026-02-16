package domain

type PolicyVerification struct {
	ProductName        string `json:"productName"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	Message            string `json:"message"`
	RegistrationNumber string `json:"registrationNumber"`
	Status             string `json:"status"`
	Success            bool   `json:"success"`
}
