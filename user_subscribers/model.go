package user_subscribers

import (
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/google/uuid"
)

type UserSubscriberAPIConfig struct {
	common.APIConfig
}

type UserSubscriber struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UserID        uuid.UUID `json:"user_id"`
	SubscriberID uuid.UUID `json:"subscriber_id"`
}