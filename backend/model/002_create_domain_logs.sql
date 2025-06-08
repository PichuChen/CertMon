-- +goose Up
CREATE TABLE IF NOT EXISTS domain_logs (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    checked_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(32) NOT NULL,
    valid_from TIMESTAMP,
    valid_to TIMESTAMP,
    days_left INTEGER,
    extra JSONB
);

-- +goose Down
DROP TABLE IF EXISTS domain_logs;
