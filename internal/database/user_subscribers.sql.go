// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_subscribers.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUserSubscriber = `-- name: CreateUserSubscriber :one
INSERT INTO user_subscribers (id, created_at, updated_at, user_id, subscriber_id)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING id, created_at, updated_at, user_id, subscriber_id
`

type CreateUserSubscriberParams struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	SubscriberID uuid.UUID
}

func (q *Queries) CreateUserSubscriber(ctx context.Context, arg CreateUserSubscriberParams) (UserSubscriber, error) {
	row := q.db.QueryRowContext(ctx, createUserSubscriber, arg.ID, arg.UserID, arg.SubscriberID)
	var i UserSubscriber
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.SubscriberID,
	)
	return i, err
}

const deleteUserSubscriber = `-- name: DeleteUserSubscriber :exec
DELETE FROM user_subscribers WHERE subscriber_id = $1 AND user_id = $2
`

type DeleteUserSubscriberParams struct {
	SubscriberID uuid.UUID
	UserID       uuid.UUID
}

func (q *Queries) DeleteUserSubscriber(ctx context.Context, arg DeleteUserSubscriberParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserSubscriber, arg.SubscriberID, arg.UserID)
	return err
}

const getUserSubscribers = `-- name: GetUserSubscribers :many
SELECT id, created_at, updated_at, user_id, subscriber_id FROM user_subscribers WHERE user_id = $1
`

func (q *Queries) GetUserSubscribers(ctx context.Context, userID uuid.UUID) ([]UserSubscriber, error) {
	rows, err := q.db.QueryContext(ctx, getUserSubscribers, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserSubscriber
	for rows.Next() {
		var i UserSubscriber
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.SubscriberID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSubscriptions = `-- name: GetUserSubscriptions :many
SELECT id, created_at, updated_at, user_id, subscriber_id FROM user_subscribers WHERE subscriber_id = $1
`

func (q *Queries) GetUserSubscriptions(ctx context.Context, subscriberID uuid.UUID) ([]UserSubscriber, error) {
	rows, err := q.db.QueryContext(ctx, getUserSubscriptions, subscriberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserSubscriber
	for rows.Next() {
		var i UserSubscriber
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.SubscriberID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
