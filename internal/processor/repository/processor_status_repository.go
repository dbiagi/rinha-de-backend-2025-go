package repository

import (
	"database/sql"
	"rinha2025/internal/domain"
)

type PaymentProcessorRepository struct {
	db *sql.DB
}

func NewProcessorStatusRepository(db *sql.DB) PaymentProcessorRepository {
	return PaymentProcessorRepository{
		db: db,
	}
}

func (r *PaymentProcessorRepository) FindStatus(tp domain.PaymentProcessorType) ([]domain.PaymentProcesor, error) {
	rows, err := r.db.Query("SELECT failing, min_response_time, code FROM payment_processor")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var processors []domain.PaymentProcesor
	for rows.Next() {
		var st domain.PaymentProcesor
		if err := rows.Scan(&st.Failing, &st.MinResponseTime, &st.Code); err != nil {
			return nil, err
		}
		processors = append(processors, st)
	}

	return processors, nil
}
