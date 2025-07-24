package service

import (
	"fmt"
	"log/slog"
	"rinha2025/internal/domain"
	"rinha2025/internal/payment/repository"
	"rinha2025/internal/processor"
	"time"
)

type PaymentService struct {
	repository           *repository.PaymentRepository
	processorService     processor.PaymentProcessorService
	maxRetry             int
	retryBackoffDuration int
}

func NewPaymentService(r *repository.PaymentRepository, ps processor.PaymentProcessorService, maxRetry int, retryBackoff int) PaymentService {
	return PaymentService{
		repository:           r,
		processorService:     ps,
		maxRetry:             maxRetry,
		retryBackoffDuration: retryBackoff,
	}
}

func (s *PaymentService) Create(request domain.PaymentCreationRequest, retryCount int) error {
	request.RequestedAt = time.Now()

	slog.Info(fmt.Sprintf("Creating payment for request=%+v", request))

	result, err := s.processorService.CreatePayment(request)

	if err != nil {
		go s.retry(request, retryCount+1)
		return err
	}

	p := domain.Payment{
		CorrelationID: request.CorrelationID,
		Amount:        request.Amount,
		RequestedAt:   request.RequestedAt,
		ProcessorID:   result.ProcessorID,
	}

	err = s.repository.Create(p)

	if err != nil {
		slog.Error("Error persisting the payment on the database", slog.String("error", err.Error()))
		return err
	}

	if retryCount != 0 {
		slog.Info(fmt.Sprintf("Payment created after %d retries. Request=%+v.", retryCount, request))
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

func (s *PaymentService) retry(pcr domain.PaymentCreationRequest, retryCount int) {
	if retryCount > s.maxRetry {
		slog.Error(fmt.Sprintf("Max retries reached for request %+v", pcr))
		return
	}

	interval := s.retryBackoffDuration * int(time.Millisecond)
	ticker := time.NewTicker(time.Duration(interval))
	defer ticker.Stop()
	<-ticker.C

	slog.Info(fmt.Sprintf("Retry number %d for request %+v", retryCount, pcr))

	s.Create(pcr, retryCount)
}
