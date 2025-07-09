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

type HealthCheckResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

type ProcessorStatus struct {
	Failing         bool
	MinResponseTime int
	Code            PaymentProcessorType
}

type PaymentProcessorType string

const (
	DefaultPaymentProcessor  PaymentProcessorType = "DEFAULT"
	FallbackPaymentProcessor PaymentProcessorType = "FALLBACK"
)
