-- +goose Up
-- +goose StatementBegin

CREATE TABLE users
(
    id                 SERIAL PRIMARY KEY,
    login_name         VARCHAR(100) NOT NULL UNIQUE,
    password_hash      TEXT         NOT NULL,
    google_auth_secret VARCHAR(64),
    email              VARCHAR(255) NOT NULL UNIQUE,
    is_verified        BOOLEAN   DEFAULT FALSE,
    scope              TEXT[] DEFAULT ARRAY['read:all', 'write:all']::TEXT[],
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users;
DROP FUNCTION update_updated_at_column;

-- +goose StatementEnd
