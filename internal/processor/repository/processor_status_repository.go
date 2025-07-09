package repository

import (
	"database/sql"
	"rinha2025/internal/processor/domain"
)

type ProcessorStatusRepository struct {
	*sql.DB
}

func NewProcessorStatusRepository(db *sql.DB) ProcessorStatusRepository {
	return ProcessorStatusRepository{
		DB: db,
	}
}

func (r *ProcessorStatusRepository) FindStatus(tp domain.PaymentProcessorType) (*domain.ProcessorStatus, error) {
	row := r.DB.QueryRow("SELECT failing, min_response_time, code FROM processor_status WHERE code = $1", string(tp))

	var st domain.ProcessorStatus
	err := row.Scan(&st.Failing, &st.MinResponseTime, &st.Code)

	if err != nil {
		return nil, err
	}

	return &st, nil
}
