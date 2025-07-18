package books

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
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

	common.JSONResponse(writer, http.StatusCreated, DatabaseBookToBookJSON(newBook))
}