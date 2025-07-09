package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"rinha2025/internal/processor/domain"
	processorerrors "rinha2025/internal/processor/errors"
)

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

func (c *PaymentProcessorClient) RequestCreatePayment(p domain.PaymentCreationRequest, ppt PaymentProcessorType) error {
	host := c.defaultHost
	if ppt == FallbackPaymentProcessor {
		host = c.fallbackHost
	}

	endpoint := fmt.Sprintf("http://%s/payments", host)
	r, err := createRequest(http.MethodPost, endpoint, p)

	if err != nil {
		slog.Error("Error creating the request", slog.String("error", err.Error()))
		return processorerrors.ErrCreatingRequest
	}

	resp, err := http.DefaultClient.Do(r)

	// TODO: Handle the diferent types of status code 4xx and 5xx
	if resp.StatusCode != http.StatusOK {
		return processorerrors.ErrUnknown
	}

	if err != nil {
		slog.Error("Processor returned an error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func createRequest(method string, url string, body any) (*http.Request, error) {
	var payload []byte = nil
	var err error

	if body != nil {
		payload, err = json.Marshal(body)
	}

	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")

	return r, nil
}
