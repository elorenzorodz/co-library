-- name: IssueBook :one
INSERT INTO book_borrows (id, issued_at, created_at, updated_at, book_id, borrower_id)
VALUES ($1, NOW(), NOW(), NOW(), $2, $3)
RETURNING id, issued_at, returned_at, created_at, updated_at, book_id, borrower_id;

-- name: GetBookBorrow :one
SELECT * FROM book_borrows WHERE book_id = $1 AND returned_at IS NULL;

-- name: ReturnBook :one
UPDATE book_borrows 
SET returned_at = NOW(), updated_at = NOW()
WHERE id = $1 AND borrower_id = $2 AND returned_at IS NULL
RETURNING id, issued_at, returned_at, created_at, updated_at, book_id, borrower_id;