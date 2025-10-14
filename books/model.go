package books

import (
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/google/uuid"
)

type BookAPIConfig struct {
	common.APIConfig
}

type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uuid.UUID `json:"user_id"`
}

type UpsertBookParameters struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}