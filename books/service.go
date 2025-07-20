package books

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/elorenzorodz/co-library/users"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (bookAPIConfig *BookAPIConfig) CreateBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	type parameters struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	params := parameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", decoderError))

		return
	}

	createBookParams := database.CreateBookParams{
		ID:        uuid.New(),
		Title:     params.Title,
		Author:    params.Author,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: userId,
	}

	newBook, createBookError := bookAPIConfig.DB.CreateBook(request.Context(), createBookParams)

	if createBookError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error creating book: %s", createBookError))

		return
	}

	// Alert subscribers about the new book.
	subscribers, getSubscribersErrors := bookAPIConfig.DB.GetUsersBySubscriberID(request.Context(), userId)

	if getSubscribersErrors != nil {
		log.Printf("Failed to get subscribers for new book alert: %s", getSubscribersErrors)
	} else {
		senderUser, getUserError := bookAPIConfig.DB.GetUserByID(request.Context(), userId)

		if getUserError != nil {
			log.Printf("Failed to get book owner details: %s", getUserError)
		} else {
			go users.DispatchNewBookAlertsSync(params.Title, subscribers, senderUser)
		}
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseBookToBookJSON(newBook))
}

func (bookAPIConfig *BookAPIConfig) GetBooks(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	getBooks, getBooksError := bookAPIConfig.DB.GetBooks(request.Context(), userId)

	if getBooksError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error getting books: %s", getBooksError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBooksToBooksJSON(getBooks))
}

func (bookAPIConfig *BookAPIConfig) GetBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["id"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid book id")

		return
	}

	getBook, getBookError := bookAPIConfig.DB.GetBook(request.Context(), bookId)

	if getBookError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error getting book: %s", getBookError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBookToBookJSON(getBook))
}

func (bookAPIConfig *BookAPIConfig) UpdateBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["id"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid book id")

		return
	}

	type parameters struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	params := parameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&params)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", decoderError))

		return
	}

	updateBookParams := database.UpdateBookParams{
		ID: bookId,
		Title: params.Title,
		Author: params.Author,
	}

	updateBook, updateBookError := bookAPIConfig.DB.UpdateBook(request.Context(), updateBookParams)

	if updateBookError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error updating book: %s", updateBookError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBookToBookJSON(updateBook))
}

func (bookAPIConfig *BookAPIConfig) DeleteBook(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	bookId, parseBookIdError := uuid.Parse(vars["id"])

	if parseBookIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid book id")

		return
	}

	deleteBookParams := database.DeleteBookParams{
		ID: bookId,
		UserID: userId,
	}

	deleteBookError := bookAPIConfig.DB.DeleteBook(request.Context(), deleteBookParams)

	if deleteBookError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error deleting book: %s", deleteBookError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, "Book successfully deleted")
}

func (bookAPIConfig *BookAPIConfig) BrowseBooks(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	browseBooks, getBooksError := bookAPIConfig.DB.BrowseBooks(request.Context())

	if getBooksError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error getting books: %s", getBooksError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBooksToBooksJSON(browseBooks))
}

func (bookAPIConfig *BookAPIConfig) BrowseBooksByUserID(writer http.ResponseWriter, request *http.Request, userId uuid.UUID) {
	vars := mux.Vars(request)
	userId, parseUserIdError := uuid.Parse(vars["id"])

	if parseUserIdError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Invalid user id")

		return
	}

	browseBooks, getBooksError := bookAPIConfig.DB.GetBooks(request.Context(), userId)

	if getBooksError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error getting books: %s", getBooksError))

		return
	}

	common.JSONResponse(writer, http.StatusOK, DatabaseBooksToBooksJSON(browseBooks))
}