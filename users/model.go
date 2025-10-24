package users

import (
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/google/uuid"
)

type UserAPIConfig struct {
	common.APIConfig
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateUserParameters struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthorized struct {
	Email string `json:"email"`
	Token string `json:"token"`
}