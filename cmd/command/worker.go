package command

import (
	"fmt"
	"log/slog"
	"rinha2025/internal/config"
	"rinha2025/internal/database"
	"rinha2025/internal/domain"
	"rinha2025/internal/processor"
	"rinha2025/internal/processor/repository"
	"time"

	"github.com/spf13/cobra"
)

const databaseAppName = "rinha-worker"

type worker struct {
	client *processor.PaymentProcessorClient
	repo   *repository.PaymentProcessorRepository
	cfg    *config.Configuration
}

func NewWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "worker",
		Run: work,
	}
}

func work(cmd *cobra.Command, args []string) {
	cfg := config.LoadConfig("worker")
	db, err := database.NewDatabase(cfg.DatabaseConfig, databaseAppName)

	if err != nil {
		slog.Error("Error connecting to database", slog.String("error", err.Error()))
		panic(1)
	}

	c := processor.NewPaymentProcessorClient()
	r := repository.NewProcessorStatusRepository(db.Connection)

	w := worker{
		client: c,
		repo:   r,
		cfg:    &cfg,
	}

	slog.Info("Worker started 🚀")
	w.start()
}

func (w *worker) start() {
	pp := w.processors()
	interval := 5 * time.Second
	ticker := time.NewTicker(interval)
	done := make(chan bool)
	for {
		select {
		case <-ticker.C:
			for _, p := range pp {
				go w.updateHealth(p)
			}
		case <-done:
			return
		}
	}
}

func (w *worker) updateHealth(p *domain.PaymentProcessor) {
	slog.Info("Checking processor health.")

	r, err := w.client.HealthCheck(*p)

	if err != nil {
		slog.Error("Error updating the health check status", slog.String("error", err.Error()))
		return
	}

	if p.Failing == r.Failing || p.MinResponseTime == r.MinResponseTime {
		return
	}

	p.Failing = r.Failing
	p.MinResponseTime = r.MinResponseTime

	slog.Info(fmt.Sprintf("Updating processor %s to Failing=%v, MinResponseTime=%d", p.Code, p.Failing, p.MinResponseTime))

	err = w.repo.UpdateHealth(p)

	if err != nil {
		slog.Error("Error updating the health status on database", slog.String("error", err.Error()))
	}
}

func (w *worker) processors() []*domain.PaymentProcessor {
	p, err := w.repo.Processors()

	if err != nil {
		panic(1)
	}

	for _, pp := range p {
		switch pp.Code {
		case domain.DefaultPaymentProcessor:
			pp.Health = domain.ProcessorHealth{
				Failing:         false,
				MinResponseTime: 0,
			}
			pp.Host = w.cfg.ProcessorConfig.DefaultHost
		case domain.FallbackPaymentProcessor:
			pp.Health = domain.ProcessorHealth{
				Failing:         false,
				MinResponseTime: 0,
			}
			pp.Host = w.cfg.ProcessorConfig.FallbackHost
		}
	}

	return p
}
