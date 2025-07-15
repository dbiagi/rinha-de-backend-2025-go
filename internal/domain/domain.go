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

type PaymentProcessorResult struct {
	Fee           float32
	ProcessorUsed PaymentProcessorType
}

type Payment struct {
	CorrelationID uuid.UUID `json:"correlationId"`
	Amount        float32   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

type HealthCheckResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

type PaymentProcesor struct {
	Failing         bool
	MinResponseTime int
	Code            PaymentProcessorType
	CheckedAt       time.Time
}

type PaymentProcessorType string

type PaymentSummary struct {
	DefaultProcessor  ProcessorSummary `json:"default"`
	FallbackProcessor ProcessorSummary `json:"fallback"`
}

type ProcessorSummary struct {
	TotalRequests float32 `json:"totalRequests"`
	TotalAmount   float32 `json:"totalAmount"`
}

const (
	DefaultPaymentProcessor  PaymentProcessorType = "DEFAULT"
	FallbackPaymentProcessor PaymentProcessorType = "FALLBACK"
)
