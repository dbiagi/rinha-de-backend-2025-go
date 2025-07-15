package service

import (
	"log/slog"
	"rinha2025/internal/domain"
	"rinha2025/internal/payment/repository"
	"rinha2025/internal/processor"
	"time"
)

type PaymentService struct {
	repository       *repository.PaymentRepository
	processorService processor.PaymentProcessorService
}

func NewPaymentService(r *repository.PaymentRepository, ps processor.PaymentProcessorService) PaymentService {
	return PaymentService{
		repository:       r,
		processorService: ps,
	}
}

func (s *PaymentService) Create(request domain.PaymentCreationRequest) error {
	request.RequestedAt = time.Now()

	err := s.processorService.CreatePayment(request)

	if err != nil {
		return err
	}

	p := domain.Payment{
		CorrelationID: request.CorrelationID,
		Amount:        request.Amount,
		RequestedAt:   request.RequestedAt,
	}

	err = s.repository.Create(p)

	if err != nil {
		slog.Error("Error persisting the payment on the database", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *PaymentService) Summary(from time.Time, to time.Time) *domain.PaymentSummary {
	summary, err := s.repository.Summary(from, to)

	if err != nil {
		slog.Error("Error fetching summary", slog.String("error", err.Error()))
		return &domain.PaymentSummary{}
	}

	return summary
}
