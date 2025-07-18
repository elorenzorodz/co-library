-- +goose Up

CREATE TABLE book_borrows (
    id UUID PRIMARY KEY,
    issued_at TIMESTAMP NOT NULL,
    returned_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    borrower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE book_borrows;
