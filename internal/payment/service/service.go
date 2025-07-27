package service

import (
	"fmt"
	"log/slog"
	"rinha2025/internal/config"
	"rinha2025/internal/domain"
	"rinha2025/internal/payment/repository"
	processorservice "rinha2025/internal/processor/service"
	"time"
)

type PaymentService struct {
	repository           *repository.PaymentRepository
	processorService     processorservice.PaymentProcessorService
	maxRetry             int
	retryBackoffDuration int
	retryQueue           chan *domain.PaymentCreationRequest
}

func NewPaymentService(r *repository.PaymentRepository, ps processorservice.PaymentProcessorService, cfg config.Configuration) PaymentService {
	svc := PaymentService{
		repository:           r,
		processorService:     ps,
		maxRetry:             cfg.MaxPaymentCreationRetries,
		retryBackoffDuration: cfg.RetryBackoffDuration,
		retryQueue:           make(chan *domain.PaymentCreationRequest, 3000),
	}

	for range cfg.RetryWorkersCount {
		go svc.startRetryQueueConsumer()
	}

	return svc
}

func (s *PaymentService) Create(request *domain.PaymentCreationRequest, retryCount int) error {
	result, err := s.processorService.CreatePayment(request)

	if err != nil {
		s.retryQueue <- request
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

func (s *PaymentService) SummaryWithFilters(from time.Time, to time.Time) *domain.PaymentSummary {
	summary, err := s.repository.SummaryWithFilters(from, to)

	if err != nil {
		slog.Error("Error fetching summary", slog.String("error", err.Error()))
		return &domain.PaymentSummary{}
	}

	return summary
}

func (s *PaymentService) Summary() *domain.PaymentSummary {
	summary, err := s.repository.Summary()

	if err != nil {
		slog.Error("Error fetching summary", slog.String("error", err.Error()))
		return &domain.PaymentSummary{}
	}

	return summary
}

func (s *PaymentService) startRetryQueueConsumer() {
	for request := range s.retryQueue {
		time.Sleep(time.Millisecond * 100)
		s.Create(request, 0)
	}
}
