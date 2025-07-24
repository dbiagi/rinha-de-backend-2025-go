package processor

import (
	"log/slog"
	"rinha2025/internal/domain"
	processorerrors "rinha2025/internal/processor/errors"
	"rinha2025/internal/processor/repository"
	"time"
)

const MaxAcceptableResponseTime = 100

type PaymentProcessorService struct {
	client            *PaymentProcessorClient
	repository        *repository.PaymentProcessorRepository
	defaultProcessor  *domain.PaymentProcessor
	fallbackProcessor *domain.PaymentProcessor
}

type PaymentCreateResult struct {
	ProcessorID int
}

type ServiceConfig struct {
	*PaymentProcessorClient
	*repository.PaymentProcessorRepository
	DefaultHost  string
	FallbackHost string
}

func NewPaymentProcessorService(c ServiceConfig) PaymentProcessorService {
	pp, err := c.PaymentProcessorRepository.Processors()

	if err != nil {
		slog.Error("Error creating the payment processors", slog.String("error", err.Error()))
		panic(1)
	}

	pps := PaymentProcessorService{
		client:     c.PaymentProcessorClient,
		repository: c.PaymentProcessorRepository,
	}

	for _, p := range pp {
		switch p.Code {
		case domain.DefaultPaymentProcessor:
			pps.defaultProcessor = p
			pps.defaultProcessor.Health = domain.ProcessorHealth{
				Failing:         false,
				MinResponseTime: 0,
			}
			pps.defaultProcessor.Host = c.DefaultHost
		case domain.FallbackPaymentProcessor:
			pps.fallbackProcessor = p
			pps.fallbackProcessor.Health = domain.ProcessorHealth{
				Failing:         false,
				MinResponseTime: 0,
			}
			pps.fallbackProcessor.Host = c.FallbackHost
		}
	}

	pps.startHealthCheckWorker()

	return pps
}

func (s *PaymentProcessorService) CreatePayment(r domain.PaymentCreationRequest) (*PaymentCreateResult, error) {
	p := s.selectProcessor()

	err := s.client.RequestCreatePayment(r, p)

	if err != nil && p.Code == domain.DefaultPaymentProcessor {
		slog.Error("Error with the default processor, using fallback", slog.String("error", err.Error()))
		errf := s.client.RequestCreatePayment(r, s.fallbackProcessor)

		if errf != nil {
			slog.Error("Error with the fallback processor", slog.String("error", errf.Error()))
			return nil, processorerrors.ErrFallbackError
		}
	}

	if err != nil && p.Code == domain.FallbackPaymentProcessor {
		slog.Error("Error with the fallback processor", slog.String("error", err.Error()))
		return nil, processorerrors.ErrFallbackError
	}

	pcr := &PaymentCreateResult{
		ProcessorID: p.ID,
	}

	return pcr, nil
}

func (s *PaymentProcessorService) selectProcessor() *domain.PaymentProcessor {
	if !s.defaultProcessor.Failing || s.defaultProcessor.MinResponseTime < MaxAcceptableResponseTime {
		return s.defaultProcessor
	}
	return s.fallbackProcessor
}

func (s *PaymentProcessorService) startHealthCheckWorker() {
	interval := time.Second * 1
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			s.checkHealth()
		}
	}()

}

func (s *PaymentProcessorService) checkHealth() {
	r, err := s.repository.Health()

	if err != nil {
		slog.Error("Error updating the health check status", slog.String("error", err.Error()))
		return
	}

	for _, h := range *r {
		switch h.ProcessorID {
		case s.defaultProcessor.ID:
			s.defaultProcessor.Failing = h.Failing
			s.defaultProcessor.MinResponseTime = h.MinResponseTime
		case s.fallbackProcessor.ID:
			s.fallbackProcessor.Failing = h.Failing
			s.fallbackProcessor.MinResponseTime = h.MinResponseTime
		}
	}
}
