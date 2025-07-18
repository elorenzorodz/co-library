-- name: CreateBook :one
INSERT INTO books (id, title, author, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, title, author, created_at, updated_at, user_id;

-- name: GetBooks :many
SELECT * FROM books WHERE user_id = $1;

-- name: GetBook :one
SELECT * FROM books WHERE user_id = $1 AND id = $2;

-- name: UpdateBook :one
UPDATE books 
SET title = $1, author = $2, updated_at = NOW() 
WHERE id = $3
RETURNING id, title, author, created_at, updated_at, user_id;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = $1 AND user_id = $2;