package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentCreationRequest struct {
	CorrelationID uuid.UUID `json:"correlationId"`
	Amount        float32   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

type PaymentCreationResponse struct {
	Message string `json:"message"`
}
