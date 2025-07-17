package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
)

func (userAPIConfig *UserAPIConfig) CreateUser(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	params := parameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error pasring JSON: %s", decoderError))

		return
	}

	hashedPassword, hashPasswordError := HashPassword(params.Password)

	if hashPasswordError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error creating user: %s", hashPasswordError))

		return
	}

	createUserParams := database.CreateUserParams {
		ID: uuid.New(),
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
		Password: hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	newUser, createUserError := userAPIConfig.DB.CreateUser(request.Context(), createUserParams)

	if createUserError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error creating user: %s", createUserError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseUserToUserJSON(newUser))
}