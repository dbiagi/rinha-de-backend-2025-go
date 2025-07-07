package handler

import (
	"net/http"

	"rinha2025/pkg/httputil"
)

type (
	HealthStatus string

	HealthCheck struct {
		Status HealthStatus `json:"status"`
	}

	HealthCheckHandler struct {
	}
)

const (
	Ok   HealthStatus = "ok"
	Down HealthStatus = "down"
)

func NewHealthCheckHandler() HealthCheckHandler {
	return HealthCheckHandler{}
}

func (h *HealthCheckHandler) Health(w http.ResponseWriter, r *http.Request) {
	hc := HealthCheck{
		Status: Ok,
	}

	httputil.NewJsonResponse(
		httputil.WithBody(hc),
	).Response(w, r)
}
