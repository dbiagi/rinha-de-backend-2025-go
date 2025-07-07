package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type (
	ContextKey string
)

const (
	TraceIdHeader                = "X-TRACE-ID"
	TraceIdContextKey ContextKey = "traceId"
)

func TraceIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := r.Header.Get(TraceIdHeader)
		if traceId == "" {
			traceId = uuid.New().String()
		}
		w.Header().Set(TraceIdHeader, traceId)

		ctx := r.Context()
		ctx = context.WithValue(ctx, TraceIdContextKey, traceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
