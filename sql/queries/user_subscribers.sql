-- name: CreateUserSubscriber :one
INSERT INTO user_subscribers (id, created_at, updated_at, user_id, subscriber_id)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING id, created_at, updated_at, user_id, subscriber_id;