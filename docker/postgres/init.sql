CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE payment_processor (
    id SERIAL PRIMARY KEY,
    failing BOOLEAN,
    min_response_time INT,
    code VARCHAR(50)
);
CREATE INDEX idx_payment_processor_code ON payment_processor (code);

INSERT INTO payment_processor (failing, min_response_time, code) VALUES
(false, 0, 'DEFAULT'),
(false, 0, 'FALLBACK');

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    correlation_id UUID,
    amount DECIMAL(10,2),
    processor_id INT REFERENCES payment_processor(id),
    requested_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_payments_processor_id ON payments (processor_id, requested_at);
