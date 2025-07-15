package processor

import (
	"log/slog"
	"rinha2025/internal/domain"
	processorerrors "rinha2025/internal/processor/errors"
	"rinha2025/internal/processor/repository"
)

const MaxAcceptableResponseTime = 100

type PaymentProcessorService struct {
	client            PaymentProcessorClient
	repository        repository.PaymentProcessorRepository
	defaultProcessor  domain.PaymentProcesor
	fallbackProcessor domain.PaymentProcesor
}

func NewPaymentProcessorService(c PaymentProcessorClient, r repository.PaymentProcessorRepository) PaymentProcessorService {
	return PaymentProcessorService{
		client:     c,
		repository: r,
		defaultProcessor: domain.PaymentProcesor{
			Failing:         false,
			Code:            domain.DefaultPaymentProcessor,
			MinResponseTime: 0,
		},
		fallbackProcessor: domain.PaymentProcesor{
			Failing:         false,
			Code:            domain.FallbackPaymentProcessor,
			MinResponseTime: 0,
		},
	}
}

func (s *PaymentProcessorService) CreatePayment(r domain.PaymentCreationRequest) error {
	s.updateHealthStatus()

	processorType := s.decideSelectedProcessor()

	err := s.client.RequestCreatePayment(r, processorType)

	if err != nil && processorType == domain.DefaultPaymentProcessor {
		slog.Error("Error with the default processor, using fallback", slog.String("error", err.Error()))
		errf := s.client.RequestCreatePayment(r, domain.FallbackPaymentProcessor)

		if errf != nil {
			slog.Error("Error with the fallback processor", slog.String("error", errf.Error()))
			return processorerrors.ErrFallbackError
		}
	}

	if err != nil && processorType == domain.FallbackPaymentProcessor {
		slog.Error("Error with the fallback processor", slog.String("error", err.Error()))
		return processorerrors.ErrFallbackError
	}

	return nil
}

func (s *PaymentProcessorService) decideSelectedProcessor() domain.PaymentProcessorType {
	if !s.defaultProcessor.Failing || s.defaultProcessor.MinResponseTime < MaxAcceptableResponseTime {
		return domain.DefaultPaymentProcessor
	}
	return domain.FallbackPaymentProcessor
}

func (s *PaymentProcessorService) updateHealthStatus() {
	processors, err := s.repository.FindStatus(domain.DefaultPaymentProcessor)

	if err != nil {
		slog.Error("Error fetching the processor status")
		return
	}

	for _, p := range processors {
		switch p.Code {
		case domain.DefaultPaymentProcessor:
			s.defaultProcessor.Failing = p.Failing
			s.defaultProcessor.MinResponseTime = p.MinResponseTime
		case domain.FallbackPaymentProcessor:
			s.fallbackProcessor.Failing = p.Failing
			s.fallbackProcessor.MinResponseTime = p.MinResponseTime
		}
	}
}
