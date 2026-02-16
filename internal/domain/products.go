package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	ProductId   int       `json:"product_id"`
	Name        string    `json:"name"`
	ProductCode string    `json:"product_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
