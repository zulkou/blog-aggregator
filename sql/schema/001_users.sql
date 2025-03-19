-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
