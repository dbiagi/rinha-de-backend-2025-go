package repository

import (
	"database/sql"
	"rinha2025/internal/domain"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment domain.Payment) error {
	query := "INSERT INTO payments (correlation_id, amount, requested_at, processor_id) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(
		query,
		payment.CorrelationID.String(),
		payment.Amount,
		payment.RequestedAt.Format(time.RFC3339),
		payment.ProcessorID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) Summary(from time.Time, to time.Time) (*domain.PaymentSummary, error) {
	query := `SELECT pp.code, SUM(p.amount), COUNT(1)
			FROM payments p
			INNER JOIN payment_processor pp ON pp.id = p.processor_id
			WHERE p.requested_at BETWEEN $1 AND $2 
			GROUP BY pp.code`
	rows, err := r.db.Query(query, from.Format(time.RFC3339), to.Format(time.RFC3339))

	if err != nil {
		return nil, err
	}

	s := &domain.PaymentSummary{
		DefaultProcessor: domain.ProcessorSummary{
			TotalRequests: 0,
			TotalAmount:   0,
		},
		FallbackProcessor: domain.ProcessorSummary{
			TotalRequests: 0,
			TotalAmount:   0,
		},
	}

	for rows.Next() {
		var code string
		var totalRequests int
		var totalAmount float32

		if err = rows.Scan(&code, &totalAmount, &totalRequests); err != nil {
			return nil, err
		}

		if code == string(domain.DefaultPaymentProcessor) {
			s.DefaultProcessor.TotalAmount = totalAmount
			s.DefaultProcessor.TotalRequests = float32(totalRequests)
		} else {
			s.FallbackProcessor.TotalAmount = totalAmount
			s.FallbackProcessor.TotalRequests = float32(totalRequests)
		}
	}

	return s, nil
}
