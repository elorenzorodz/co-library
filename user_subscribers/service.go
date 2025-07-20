package user_subscribers

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) CreateUserSubscriber(writer http.ResponseWriter, request *http.Request, subscriberId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["user_id"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid user id")

		return
	}

	createUserSubscriberParam := database.CreateUserSubscriberParams {
		ID: uuid.New(),
		UserID: userId,
		SubscriberID: subscriberId,
	}

	newUserSubscriber, createUserSubscriberError := userSubscriberAPIConfig.DB.CreateUserSubscriber(request.Context(), createUserSubscriberParam)

	if createUserSubscriberError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error subscribiing to user: %s", createUserSubscriberError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseUserSubscriberToUserSubscriberJSON(newUserSubscriber))
}

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) DeleteUserSubscriber(writer http.ResponseWriter, request *http.Request, subscriberId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["user_id"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid user id")

		return
	}

	deleteUserSubscriberParam := database.DeleteUserSubscriberParams {
		SubscriberID: subscriberId,
		UserID: userId,
	}

	deleteUserSubscriberError := userSubscriberAPIConfig.DB.DeleteUserSubscriber(request.Context(), deleteUserSubscriberParam)

	if deleteUserSubscriberError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error removing subscription to user: %s", deleteUserSubscriberError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, "User subscription successfully deleted")
}

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) GetUserSubscribers(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	userSubscribers, getUserSubscribersError := userSubscriberAPIConfig.DB.GetUserSubscribers(request.Context(), userId)

	if getUserSubscribersError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error getting subscribers: %s", getUserSubscribersError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseUserSubscribersToUserSubscribersJSON(userSubscribers))
}