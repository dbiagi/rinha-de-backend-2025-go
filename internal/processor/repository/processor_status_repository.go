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

func (r *PaymentProcessorRepository) Processors() ([]*domain.PaymentProcessor, error) {
	rows, err := r.db.Query("SELECT id, failing, min_response_time, code FROM payment_processor")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var processors []*domain.PaymentProcessor
	for rows.Next() {
		var st domain.PaymentProcessor
		if err := rows.Scan(&st.ID, &st.Failing, &st.MinResponseTime, &st.Code); err != nil {
			return nil, err
		}
		processors = append(processors, &st)
	}

	return processors, nil
}

func (r *PaymentProcessorRepository) UpdateHealth(p *domain.PaymentProcessor) error {
	query := `
		UPDATE payment_processor
		SET failing = $2, min_response_time = $3
		WHERE id = $1
	`
	_, err := r.db.Exec(
		query,
		p.ID,
		p.Failing,
		p.MinResponseTime,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentProcessorRepository) Health() (*[]domain.ProcessorHealth, error) {
	query := `
		SELECT failing, min_response_time, id
		FROM payment_processor
	`
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, nil
	}

	result := []domain.ProcessorHealth{}
	for rows.Next() {
		ph := domain.ProcessorHealth{}
		if err := rows.Scan(&ph.Failing, &ph.MinResponseTime, &ph.ProcessorID); err != nil {
			return nil, err
		}
		result = append(result, ph)
	}

	return &result, nil
}
