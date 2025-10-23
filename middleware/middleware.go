package middleware

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/google/uuid"
)

type AuthHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func Authorization(apiConfig *common.APIConfig, handler AuthHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		jwt, jwtError := common.GetJWT(request.Header)

		if jwtError != nil {
			common.ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("authentication error: %s", jwtError))

			return
		}

		email, extractEmailClaimError := common.ValidateJWTAndGetEmailClaim(jwt, apiConfig.JWTValidationKey)

		if extractEmailClaimError != nil {
			common.ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("authentication error: %s", extractEmailClaimError))

			return
		}

		getUser, getUserError := apiConfig.DB.GetUserByEmail(request.Context(), email)

		if getUserError != nil {
			common.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("authentication error: %s", getUserError))

			return
		}

		handler(writer, request, getUser.ID)
	}
}