package book_borrows

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (bookBorrowAPIConfig *BookBorrowAPIConfig) IssueBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["bookId"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid book id")

		return
	}

	getBook, getBookError := bookBorrowAPIConfig.DB.GetBook(request.Context(), bookId)

	if getBookError != nil {
		if getBookError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusNotFound, "book not found")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "failed to get book details, please try again in a few minutes")
		}

		return
	}

	// Borrower cannot borrow their own book.
	if userId == getBook.UserID {
		common.ErrorResponse(writer, http.StatusForbidden, "you cannot borrow your own book")

		return
	}

	_, getBookBorrowError := bookBorrowAPIConfig.DB.GetBookBorrow(request.Context(), bookId)

	if getBookBorrowError != nil && getBookBorrowError != sql.ErrNoRows {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error issuing book: %s", getBookBorrowError))

		return
	} else if getBookBorrowError == nil {
		common.ErrorResponse(writer, http.StatusConflict, "book is currently issued to another borrower")

		return
	}

	issueBookParams := database.IssueBookParams{
		ID:         uuid.New(),
		BookID:     getBook.ID,
		BorrowerID: userId,
	}

	issueBook, issueBookError := bookBorrowAPIConfig.DB.IssueBook(request.Context(), issueBookParams)

	if issueBookError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error issuing book: %s", issueBookError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseBookBorrowToBookBorrowJSON(issueBook))
}

func (bookBorrowAPIConfig *BookBorrowAPIConfig) ReturnBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookBorrowId, parseBookBorrowIdError := uuid.Parse(vars["bookBorrowId"])

	if parseBookBorrowIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid book borrow id")

		return
	}

	returnBookParams := database.ReturnBookParams {
		ID: bookBorrowId,
		BorrowerID: userId,
	}

	returnBook, returnBookError := bookBorrowAPIConfig.DB.ReturnBook(request.Context(), returnBookParams)

	if returnBookError != nil {
		if returnBookError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusBadRequest, "failed to return book: record not found, unauthorized, or already returned")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "failed to return book, please try again in a few minutes")
		}

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBookBorrowToBookBorrowJSON(returnBook))
}