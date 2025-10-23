package user_subscribers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) CreateUserSubscriber(writer http.ResponseWriter, request *http.Request, subscriberId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["userId"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid user id")

		return
	}

	if userId  == subscriberId {
		common.ErrorResponse(writer, http.StatusBadRequest, "cannot subscribe to self")

		return
	}

	// Check if the user exists to which the subscriber is trying to subscribe.
	_, getUserByIDError := userSubscriberAPIConfig.DB.GetUserByID(request.Context(), userId)

	if getUserByIDError != nil {
		if getUserByIDError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusNotFound, "the user you are trying to subscribe is not found")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "failed to subscribe to user, please try again in a few minutes")
		}

		return
	}

	// Check if user is already subscribed.
	getUserSubscriberParams := database.GetUserSubscriberParams {
		SubscriberID: subscriberId,
		UserID: userId,
	}

	_, getUserSubscriberError := userSubscriberAPIConfig.DB.GetUserSubscriber(request.Context(), getUserSubscriberParams)

	if getUserSubscriberError != nil && getUserSubscriberError != sql.ErrNoRows {
		common.ErrorResponse(writer, http.StatusInternalServerError, "failed to subscribe to user, please try again in a few minutes")

		return
	} else if getUserSubscriberError == nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "you are already subscribed to user")

		return
	}

	createUserSubscriberParam := database.CreateUserSubscriberParams{
		ID:           uuid.New(),
		UserID:       userId,
		SubscriberID: subscriberId,
	}

	newUserSubscriber, createUserSubscriberError := userSubscriberAPIConfig.DB.CreateUserSubscriber(request.Context(), createUserSubscriberParam)

	if createUserSubscriberError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error subscribiing to user: %s", createUserSubscriberError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseUserSubscriberToUserSubscriberJSON(newUserSubscriber))
}

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) DeleteUserSubscriber(writer http.ResponseWriter, request *http.Request, subscriberId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["userId"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid user id")

		return
	}

	deleteUserSubscriberParam := database.DeleteUserSubscriberParams{
		SubscriberID: subscriberId,
		UserID:       userId,
	}

	rowsAffected, deleteUserSubscriberError := userSubscriberAPIConfig.DB.DeleteUserSubscriber(request.Context(), deleteUserSubscriberParam)

	if deleteUserSubscriberError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error unsubscribing to user: %s", deleteUserSubscriberError))

		return
	}

	if rowsAffected == 0 {
		common.ErrorResponse(writer, http.StatusNotFound, "user subscription not found")

        return
	}

	common.JSONResponse(writer, http.StatusOK, "user subscription successfully deleted")
}

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) GetUserSubscribers(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	userSubscribers, getUserSubscribersError := userSubscriberAPIConfig.DB.GetUserSubscribers(request.Context(), userId)

	if getUserSubscribersError != nil {
		if getUserSubscribersError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusOK, "no user subscribers")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error getting subscribers: %s", getUserSubscribersError))
		}

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseUserSubscribersToUserSubscribersJSON(userSubscribers))
}

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) GetUserSubscriptions(writer http.ResponseWriter, request *http.Request, subscriberId uuid.UUID) {
	userSubscriptions, getUserSubscriptionsError := userSubscriberAPIConfig.DB.GetUserSubscriptions(request.Context(), subscriberId)

	if getUserSubscriptionsError != nil {
		if getUserSubscriptionsError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusOK, "no user subscriptions found")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error getting subscriptions: %s", getUserSubscriptionsError))
		}

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseUserSubscribersToUserSubscribersJSON(userSubscriptions))
}