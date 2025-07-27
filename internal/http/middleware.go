package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func LogLatencyMiddleware(id string, f handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		f(w, r)
		d := time.Since(n)
		slog.Info(fmt.Sprintf("Request id %s. Duration: %d ms", id, d.Milliseconds()))
	}
}
