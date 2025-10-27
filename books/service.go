package books

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/elorenzorodz/co-library/users"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (bookAPIConfig *BookAPIConfig) CreateBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	upsertBookParameters := UpsertBookParameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&upsertBookParameters)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %s", decoderError))

		return
	}

	if strings.TrimSpace(upsertBookParameters.Title) == "" || strings.TrimSpace(upsertBookParameters.Author) == "" {
		common.ErrorResponse(writer, http.StatusBadRequest, "title and author are required")
		
		return
	}

	createBookParams := database.CreateBookParams{
		ID:        uuid.New(),
		Title:     upsertBookParameters.Title,
		Author:    upsertBookParameters.Author,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userId,
	}

	newBook, createBookError := bookAPIConfig.DB.CreateBook(request.Context(), createBookParams)

	if createBookError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error creating book: %s", createBookError))

		return
	}

	// Alert subscribers about the new book.
	subscribers, getSubscribersErrors := bookAPIConfig.DB.GetUsersBySubscriberID(request.Context(), userId)

	if getSubscribersErrors != nil {
		log.Printf("failed to get subscribers for new book alert: %s", getSubscribersErrors)
	} else {
		senderUser, getUserError := bookAPIConfig.DB.GetUserByID(request.Context(), userId)

		if getUserError != nil {
			log.Printf("failed to get book owner details: %s", getUserError)
		} else {
			go users.DispatchNewBookAlertsSync(upsertBookParameters.Title, subscribers, senderUser, bookAPIConfig.APIConfig.MailgunAPIKey, bookAPIConfig.APIConfig.MailgunSendingDomain)
		}
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseBookToBookJSON(newBook))
}

func (bookAPIConfig *BookAPIConfig) GetBooks(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	getBooks, getBooksError := bookAPIConfig.DB.GetBooks(request.Context(), userId)

	if getBooksError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error getting books: %s", getBooksError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBooksToBooksJSON(getBooks))
}

func (bookAPIConfig *BookAPIConfig) GetBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["bookId"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid book id")

		return
	}

	getBook, getBookError := bookAPIConfig.DB.GetBook(request.Context(), bookId)

	if getBookError != nil {
		if getBookError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusNotFound, "book not found")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "error getting book details, please try again in a few minutes")
		}

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBookToBookJSON(getBook))
}

func (bookAPIConfig *BookAPIConfig) UpdateBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["bookId"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid book id")

		return
	}

	upsertBookParameters := UpsertBookParameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&upsertBookParameters)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %s", decoderError))

		return
	}

	updateBookParams := database.UpdateBookParams{
		Title:  upsertBookParameters.Title,
		Author: upsertBookParameters.Author,
		ID:     bookId,
		UserID: userId,
	}

	updateBook, updateBookError := bookAPIConfig.DB.UpdateBook(request.Context(), updateBookParams)

	if updateBookError != nil {
		if updateBookError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusNotFound, "book not found")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "error updating book details, please try again in a few minutes")
		}

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBookToBookJSON(updateBook))
}

func (bookAPIConfig *BookAPIConfig) DeleteBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["bookId"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid book id")

		return
	}

	deleteBookParams := database.DeleteBookParams{
		ID:     bookId,
		UserID: userId,
	}

	rowsAffected, deleteBookError := bookAPIConfig.DB.DeleteBook(request.Context(), deleteBookParams)

	if deleteBookError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error deleting book: %s", deleteBookError))

		return
	}

	if rowsAffected == 0 {
		common.ErrorResponse(writer, http.StatusNotFound, "book not found")

        return
	}

	common.JSONResponse(writer, http.StatusOK, "book successfully deleted")
}

func (bookAPIConfig *BookAPIConfig) BrowseBooks(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	browseBooks, getBooksError := bookAPIConfig.DB.BrowseBooks(request.Context())

	if getBooksError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error getting books: %s", getBooksError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBooksToBooksJSON(browseBooks))
}

func (bookAPIConfig *BookAPIConfig) BrowseBooksByUserID(writer http.ResponseWriter, request *http.Request, uId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["userId"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "invalid user id")

		return
	}

	browseBooks, getBooksError := bookAPIConfig.DB.GetBooks(request.Context(), userId)

	if getBooksError != nil {
		if getBooksError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusNotFound, "user not found or has no books")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "error getting books, please try again in a few minutes")
		}

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBooksToBooksJSON(browseBooks))
}