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
	"rinha2025/internal/processor"
	"rinha2025/internal/processor/repository"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"rinha2025/internal/database"
	paymenthandler "rinha2025/internal/payment/handler"
	paymentrepo "rinha2025/internal/payment/repository"
	paymentservice "rinha2025/internal/payment/service"
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
		paymenthandler.PaymentHandler
	}
)

func NewServer(c config.Configuration) *Server {
	return &Server{
		Configuration: c,
	}
}

func (s *Server) Start() {
	server, router := createServer(s.Configuration.WebConfig)
	db := connectToDatabase(s.DatabaseConfig)
	handlers := createHandlers(s.Configuration, db)
	registerRoutesAndMiddlewares(router, handlers)
	configureGracefullShutdown(server, db, s.Configuration.WebConfig)
}

func connectToDatabase(cfg config.DatabaseConfig) *database.Database {
	db, err := database.NewDatabase(cfg)

	if err != nil {
		os.Exit(2)
	}

	return db
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
	router.HandleFunc("/health", h.HealthCheckHandler.Health).Methods(http.MethodGet)
	router.HandleFunc("/payments", h.PaymentHandler.CreatePayment).Methods(http.MethodPost)
	router.HandleFunc("/payments-summary", h.PaymentHandler.Summary).Methods(http.MethodGet)
	router.Use(handlers.CompressHandler)
}

func configureGracefullShutdown(server *http.Server, db *database.Database, webConfig config.WebConfig) {
	slog.Info("Configuring graceful shutdown.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, signalsToListenTo...)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), webConfig.ShutdownTimeout)
	defer cancel()

	slog.Info("Shutting down server")
	server.Shutdown(ctx)

	slog.Info("Closing database connection")
	db.Close()

	os.Exit(0)
}

func createHandlers(c config.Configuration, db *database.Database) appHandlers {
	return appHandlers{
		HealthCheckHandler: healthhandler.NewHealthCheckHandler(),
		PaymentHandler:     createPaymentHandler(c.ProcessorConfig, db),
	}
}

func createPaymentHandler(c config.ProcessorConfig, db *database.Database) paymenthandler.PaymentHandler {
	repository := repository.NewProcessorStatusRepository(db.Connection)
	client := processor.NewPaymentProcessorClient(c.DefaultHost, c.FallbackHost)
	service := processor.NewPaymentProcessorService(client, repository)

	paymentRepo := paymentrepo.NewPaymentRepository(db.Connection)
	paymentService := paymentservice.NewPaymentService(paymentRepo, service)
	handler := paymenthandler.NewPaymentHandler(paymentService)

	return handler
}
