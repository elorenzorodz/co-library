-- +goose Up

CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE books;
