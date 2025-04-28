-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    uuid VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) DEFAULT NULL,
    age INT DEFAULT NULL,
    gender VARCHAR(6) DEFAULT NULL,
    country_id VARCHAR(2) DEFAULT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_name ON users (name);
CREATE INDEX IF NOT EXISTS idx_users_surname ON users (surname);
CREATE INDEX IF NOT EXISTS idx_users_patronymic ON users (patronymic);
CREATE INDEX IF NOT EXISTS idx_users_age ON users (age);
CREATE INDEX IF NOT EXISTS idx_users_gender ON users (gender);
CREATE INDEX IF NOT EXISTS idx_users_country_id ON users (country_id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER set_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at_trigger ON users;
DROP FUNCTION IF EXISTS set_updated_at();
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_users_surname;
DROP INDEX IF EXISTS idx_users_patronymic;
DROP INDEX IF EXISTS idx_users_age;
DROP INDEX IF EXISTS idx_users_gender;
DROP INDEX IF EXISTS idx_users_country_id;
DROP TABLE IF EXISTS users;