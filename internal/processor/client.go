package processor

import "rinha2025/internal/processor/domain"

type PaymentProcessorClient struct {
	defaultHost  string
	fallbackHost string
}

type PaymentProcessorType string

const (
	DefaultPaymentProcessor  PaymentProcessorType = "DEFAULT"
	FallbackPaymentProcessor PaymentProcessorType = "FALLBACK"
)

func NewPaymentProcessorClient(defaultHost string, fallbackHost string) PaymentProcessorClient {
	return PaymentProcessorClient{
		defaultHost:  defaultHost,
		fallbackHost: fallbackHost,
	}
}

func (c *PaymentProcessorClient) RequestCreatePayment(p domain.PaymentCreationRequest, ppt PaymentProcessorType) {

}
