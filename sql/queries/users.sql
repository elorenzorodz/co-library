-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, first_name, last_name, email, password, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsersBySubscriberID :many
SELECT u.*
FROM users AS u
LEFT JOIN user_subscribers AS us
ON us.subscriber_id = u.ID
WHERE us.user_id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;