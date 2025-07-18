-- name: CreateBook :one
INSERT INTO books (id, title, author, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, title, author, created_at, updated_at, user_id;

-- name: GetBooks :many
SELECT * FROM books WHERE user_id = $1;