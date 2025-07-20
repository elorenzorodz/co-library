package user_subscribers

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
)

func DatabaseUserSubscriberToUserSubscriberJSON(databaseUserSubscriber database.UserSubscriber) UserSubscriber {
	return UserSubscriber{
		ID:            databaseUserSubscriber.ID,
		CreatedAt:     databaseUserSubscriber.CreatedAt,
		UpdatedAt:     databaseUserSubscriber.UpdatedAt,
		UserID:        databaseUserSubscriber.UserID,
		SubscriberID: databaseUserSubscriber.SubscriberID,
	}
}

type UserSubscriberAuthHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (userSubscriberAPIConfig *UserSubscriberAPIConfig) Authorization(handler UserSubscriberAuthHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		jwt, jwtError := common.GetJWT(request.Header)

		if jwtError != nil {
			common.ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("Authentication error: %s", jwtError))

			return
		}

		email, extractEmailClaimError := common.ValidateJWTAndGetEmailClaim(jwt)

		if extractEmailClaimError != nil {
			common.ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("Authentication error: %s", extractEmailClaimError))

			return
		}

		getUser, getUserError := userSubscriberAPIConfig.DB.GetUserByEmail(request.Context(), email)

		if getUserError != nil {
			common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Authentication error: %s", getUserError))

			return
		}

		handler(writer, request, getUser.ID)
	}
}