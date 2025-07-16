package repository

import (
	"database/sql"
	"rinha2025/internal/domain"
)

type PaymentProcessorRepository struct {
	db *sql.DB
}

func NewProcessorStatusRepository(db *sql.DB) *PaymentProcessorRepository {
	return &PaymentProcessorRepository{
		db: db,
	}
}

func (r *PaymentProcessorRepository) Processors() ([]domain.PaymentProcessor, error) {
	rows, err := r.db.Query("SELECT id, failing, min_response_time, code FROM payment_processor")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var processors []domain.PaymentProcessor
	for rows.Next() {
		var st domain.PaymentProcessor
		if err := rows.Scan(&st.ID, &st.Failing, &st.MinResponseTime, &st.Code); err != nil {
			return nil, err
		}
		processors = append(processors, st)
	}

	return processors, nil
}
