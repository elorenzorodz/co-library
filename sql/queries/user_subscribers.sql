-- name: CreateUserSubscriber :one
INSERT INTO user_subscribers (id, created_at, updated_at, user_id, subscriber_id)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING id, created_at, updated_at, user_id, subscriber_id;

-- name: DeleteUserSubscriber :exec
DELETE FROM user_subscribers WHERE subscriber_id = $1 AND user_id = $2;

-- name: GetUserSubscribers :many
SELECT * FROM user_subscribers WHERE user_id = $1;