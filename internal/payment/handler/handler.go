package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"rinha2025/internal/domain"
	"rinha2025/internal/payment/service"
	"rinha2025/pkg/httputil"
	"time"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(s service.PaymentService) PaymentHandler {
	return PaymentHandler{
		service: s,
	}
}

func (ph *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var cpr domain.PaymentCreationRequest
	err := json.NewDecoder(r.Body).Decode(&cpr)

	if err != nil {
		slog.Error("Error deserializing the request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		ph.service.Create(cpr)
	}()

	w.WriteHeader(http.StatusCreated)
}

func (ph *PaymentHandler) Summary(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		slog.Error("Missing 'from' or 'to' query parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		slog.Error("Invalid 'from' query parameter", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		slog.Error("Invalid 'to' query parameter", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	summary := ph.service.Summary(fromTime, toTime)

	httputil.NewJsonResponse(httputil.WithBody(summary)).Response(w, r)
}
