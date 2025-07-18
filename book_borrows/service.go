package book_borrows

import (
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (bookBorrowAPIConfig *BookBorrowAPIConfig) IssueBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["id"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid book id")

		return
	}

	getBook, getBookError := bookBorrowAPIConfig.DB.GetBook(request.Context(), bookId)

	if getBookError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error issuing book: %s", getBookError))

		return
	}

	getBookBorrow, getBookBorrowError := bookBorrowAPIConfig.DB.GetBookBorrow(request.Context(), bookId)

	if getBookBorrowError != nil &&  getBookBorrowError.Error() != "sql: no rows in result set" {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error issuing book: %s", getBookBorrowError))

		return
		
	}

	if getBookBorrow.BookID == getBook.ID {
		common.ErrorResponse(writer, http.StatusBadRequest, "Book is not yet returned by another borrower")

		return
	}

	issueBookParams := database.IssueBookParams{
		ID:        uuid.New(),
		BookID: getBook.ID,
		BorrowerID:    userId,
	}

	issueBook, issueBookError := bookBorrowAPIConfig.DB.IssueBook(request.Context(), issueBookParams)

	if issueBookError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error issuing book: %s", issueBookError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseBookBorrowToBookBorrowJSON(issueBook))
}