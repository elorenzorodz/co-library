package books

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
)

func DatabaseBookToBookJSON(databaseBook database.Book) Book {
	return Book{
		ID:        databaseBook.ID,
		Title:     databaseBook.Title,
		Author:    databaseBook.Author,
		CreatedAt: databaseBook.CreatedAt,
		UpdatedAt: databaseBook.UpdatedAt,
		UserID:    databaseBook.UserID,
	}
}

type BookAuthHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (bookAPIConfig *BookAPIConfig) Authorization(handler BookAuthHandler) http.HandlerFunc {
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

		getUser, getUserError := bookAPIConfig.DB.GetUserByEmail(request.Context(), email)

		if getUserError != nil {
			common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Authentication error: %s", getUserError))

			return
		}

		handler(writer, request, getUser.ID)
	}
}