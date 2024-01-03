-- +goose Up
CREATE TABLE IF NOT EXISTS reset_passwords(
    id VARCHAR(36) PRIMARY KEY,
    token VARCHAR(6) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TRIGGER update_reset_passwords_updated_on BEFORE UPDATE ON reset_passwords FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
DROP TRIGGER update_reset_passwords_updated_on ON reset_passwords;
DROP TABLE IF EXISTS reset_passwords;