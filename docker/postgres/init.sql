CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE processor_status (
    id SERIAL PRIMARY KEY,
    failing BOOLEAN,
    min_response_time INT,
    code VARCHAR(50),
    status VARCHAR(50)
);
CREATE INDEX idx_processor_status_code ON processor_status (code);

INSERT INTO processor_status (failing, min_response_time, code) VALUES
(false, 0, 'DEFAULT'),
(false, 0, 'FALLBACK');
