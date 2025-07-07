package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"rinha2025/internal/config"
	healthhandler "rinha2025/internal/health/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"rinha2025/pkg/middleware"
)

var (
	signalsToListenTo = []os.Signal{
		syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
	}
)

type (
	Server struct {
		config.Configuration
		*http.Server
		*mux.Router
	}

	appHandlers struct {
		healthhandler.HealthCheckHandler
	}
)

func NewServer(c config.Configuration) *Server {
	return &Server{
		Configuration: c,
	}
}

func (s *Server) Start() {
	server, router := createServer(s.Configuration.WebConfig)

	handlers := createHandlers(s.Configuration)

	registerRoutesAndMiddlewares(router, handlers)
	configureGracefullShutdown(server, s.Configuration.WebConfig)
}

func (s *Server) ForceShutdown() {
	s.Server.Shutdown(context.Background())
}

func createServer(webConfig config.WebConfig) (*http.Server, *mux.Router) {
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", webConfig.Port),
		Handler:      router,
		IdleTimeout:  webConfig.IdleTimeout,
		ReadTimeout:  webConfig.ReadTimeout,
		WriteTimeout: webConfig.WriteTimeout,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err.Error() != "http: Server closed" {
			slog.Error("Error starting server.", slog.String("error", err.Error()))
		}
	}()

	return srv, router
}

func registerRoutesAndMiddlewares(router *mux.Router, h appHandlers) {
	router.Use(middleware.TraceIdMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/health", h.HealthCheckHandler.Health).Methods("GET")
	router.Use(handlers.CompressHandler)
}

func configureGracefullShutdown(server *http.Server, webConfig config.WebConfig) {
	slog.Info("Configuring graceful shutdown.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, signalsToListenTo...)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), webConfig.ShutdownTimeout)
	defer cancel()

	server.Shutdown(ctx)
	slog.Info("Shutting down server")
	os.Exit(0)
}

func createHandlers(c config.Configuration) appHandlers {
	return appHandlers{
		HealthCheckHandler: healthhandler.NewHealthCheckHandler(),
	}
}
