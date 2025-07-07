package httputil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type (
	JsonResponse struct {
		statusCode *int
		body       *any
	}

	Option func(*JsonResponse)
)

func NewJsonResponse(opts ...Option) *JsonResponse {
	jr := &JsonResponse{}

	for _, opt := range opts {
		opt(jr)
	}

	return jr
}

func (jr *JsonResponse) Response(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if jr.statusCode != nil {
		w.WriteHeader(*jr.statusCode)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if jr.body != nil {
		writeBoby(*jr.body, w)
	}
}

func WithStatusCode(statusCode int) Option {
	return func(jr *JsonResponse) {
		jr.statusCode = &statusCode
	}
}

func WithBody(body any) Option {
	return func(jr *JsonResponse) {
		jr.body = &body
	}
}

func writeBoby(body any, w http.ResponseWriter) {
	bytes, err := json.Marshal(body)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Error marshalling response: %v", err.Error()),
			slog.String("error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(bytes)
}
