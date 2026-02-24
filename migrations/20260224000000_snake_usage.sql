-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS snake_usage (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    counter INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_snake_usage_user_id ON snake_usage(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_snake_usage_user_id;
DROP TABLE IF EXISTS snake_usage;

-- +goose StatementEnd
