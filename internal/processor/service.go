package processor

import "rinha2025/internal/processor/repository"

type TransactionProcessorService struct {
	PaymentProcessorClient
	repository.ProcessorStatusRepository
}

func NewTransactionProcessorService(c PaymentProcessorClient, r repository.ProcessorStatusRepository) TransactionProcessorService {
	return TransactionProcessorService{
		PaymentProcessorClient:    c,
		ProcessorStatusRepository: r,
	}
}
