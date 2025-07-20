package common

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type AuthHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (apiConfig *APIConfig) Authorization(handler AuthHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		jwt, jwtError := GetJWT(request.Header)

		if jwtError != nil {
			ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("Authentication error: %s", jwtError))

			return
		}

		email, extractEmailClaimError := ValidateJWTAndGetEmailClaim(jwt)

		if extractEmailClaimError != nil {
			ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("Authentication error: %s", extractEmailClaimError))

			return
		}

		getUser, getUserError := apiConfig.DB.GetUserByEmail(request.Context(), email)

		if getUserError != nil {
			ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Authentication error: %s", getUserError))

			return
		}

		handler(writer, request, getUser.ID)
	}
}