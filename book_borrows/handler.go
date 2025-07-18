package book_borrows

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
)

func DatabaseBookBorrowToBookBorrowJSON(databaseBookBorrow database.BookBorrow) BookBorrow {
	return BookBorrow{
		ID:        databaseBookBorrow.ID,
		IssuedAt:     databaseBookBorrow.IssuedAt,
		ReturnedAt:    databaseBookBorrow.ReturnedAt,
		CreatedAt: databaseBookBorrow.CreatedAt,
		UpdatedAt: databaseBookBorrow.UpdatedAt,
		BookID:    databaseBookBorrow.BookID,
		BorrowerID:    databaseBookBorrow.BorrowerID,
	}
}

type BookBorrowAuthHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (bookBorrowAPIConfig *BookBorrowAPIConfig) Authorization(handler BookBorrowAuthHandler) http.HandlerFunc {
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

		getUser, getUserError := bookBorrowAPIConfig.DB.GetUserByEmail(request.Context(), email)

		if getUserError != nil {
			common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Authentication error: %s", getUserError))

			return
		}

		handler(writer, request, getUser.ID)
	}
}