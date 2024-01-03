-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    nick_name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    photo_url VARCHAR(200),
    password VARCHAR(200) NOT NULL,
    salt VARCHAR(200) NOT NULL,
    role VARCHAR(200) NOT NULL DEFAULT 'free',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE TRIGGER update_users_updated_on BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
DROP TRIGGER update_users_updated_on ON users;
DROP TABLE IF EXISTS users;