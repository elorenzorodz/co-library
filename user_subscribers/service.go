package user_subscribers

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) CreateUserSubscriber(writer http.ResponseWriter, request *http.Request, borrowerId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["user_id"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid user id")

		return
	}

	createUserSubscriberParam := database.CreateUserSubscriberParams {
		ID: uuid.New(),
		UserID: userId,
		SubscriberID: borrowerId,
	}

	newUserSubscriber, createUserSubscriberError := userSubscriberAPIConfig.DB.CreateUserSubscriber(request.Context(), createUserSubscriberParam)

	if createUserSubscriberError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error subscribiing to user: %s", createUserSubscriberError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseUserSubscriberToUserSubscriberJSON(newUserSubscriber))
}