CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE processor_status (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50),
    status VARCHAR(50)
);
CREATE INDEX idx_processor_status_code ON processor_status (code);
