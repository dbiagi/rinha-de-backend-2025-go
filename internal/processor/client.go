package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"rinha2025/internal/domain"
	processorerrors "rinha2025/internal/processor/errors"
)

type PaymentProcessorClient struct {
	defaultHost  string
	fallbackHost string
}

func NewPaymentProcessorClient(defaultHost string, fallbackHost string) PaymentProcessorClient {
	return PaymentProcessorClient{
		defaultHost:  defaultHost,
		fallbackHost: fallbackHost,
	}
}

func (c *PaymentProcessorClient) RequestCreatePayment(p domain.PaymentCreationRequest, ppt domain.PaymentProcessorType) error {
	host := c.defaultHost
	if ppt == domain.FallbackPaymentProcessor {
		host = c.fallbackHost
	}

	endpoint := fmt.Sprintf("http://%s/payments", host)
	r, err := createRequest(http.MethodPost, endpoint, p)

	if err != nil {
		slog.Error("Error creating the request", slog.String("error", err.Error()))
		return processorerrors.ErrCreatingRequest
	}

	resp, err := http.DefaultClient.Do(r)

	if is5xxResponse(resp.StatusCode) {
		return processorerrors.ErrCreatingRequest
	}

	if !is2xxReponse(resp.StatusCode) {
		return processorerrors.ErrUnknown
	}

	defer resp.Body.Close()

	if err != nil {
		slog.Error("Processor returned an error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (c *PaymentProcessorClient) HealthCheck() (*domain.HealthCheckResponse, error) {
	endpoint := fmt.Sprintf("http://%s/payments/service-health", c.defaultHost)

	r, err := createRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		slog.Error("Error creating the request", slog.String("error", err.Error()))
		return nil, processorerrors.ErrCreatingRequest
	}

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		slog.Error("Health check returned an error")
		return nil, processorerrors.ErrUnknown
	}

	defer resp.Body.Close()

	hc := &domain.HealthCheckResponse{}
	err = json.NewDecoder(resp.Body).Decode(hc)
	if err != nil {
		slog.Error("Error deserialing the response body", slog.String("error", err.Error()))
		return nil, processorerrors.ErrUnknown
	}

	return hc, nil
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

func is5xxResponse(status int) bool {
	return status >= 500 && status < 600
}

func is2xxReponse(status int) bool {
	return status >= 200 && status < 300
}
