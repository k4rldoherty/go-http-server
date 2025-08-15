-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE refresh_tokens;
