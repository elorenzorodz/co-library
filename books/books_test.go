package books

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MockQueries struct {
	*common.BaseMock

	CreateBookFunc func(ctx context.Context, arg database.CreateBookParams) (database.Book, error)
	GetBookFunc    func(ctx context.Context, id uuid.UUID) (database.Book, error)
	UpdateBookFunc func(ctx context.Context, arg database.UpdateBookParams) (database.Book, error)
	DeleteBookFunc func(ctx context.Context, arg database.DeleteBookParams) (int64, error)
}

func (mockQueries *MockQueries) CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
	if mockQueries.CreateBookFunc != nil {
		return mockQueries.CreateBookFunc(ctx, arg)
	}

	return mockQueries.BaseMock.CreateBook(ctx, arg)
}

func (mockQueries *MockQueries) GetBook(ctx context.Context, id uuid.UUID) (database.Book, error) {
	if mockQueries.GetBookFunc != nil {
		return mockQueries.GetBookFunc(ctx, id)
	}

	return mockQueries.BaseMock.GetBook(ctx, id)
}

func (mockQueries *MockQueries) UpdateBook(ctx context.Context, arg database.UpdateBookParams) (database.Book, error) {
	if mockQueries.UpdateBookFunc != nil {
		return mockQueries.UpdateBookFunc(ctx, arg)
	}

	return mockQueries.BaseMock.UpdateBook(ctx, arg)
}

func (mockQueries *MockQueries) DeleteBook(ctx context.Context, arg database.DeleteBookParams) (int64, error) {
	if mockQueries.DeleteBookFunc != nil {
		return mockQueries.DeleteBookFunc(ctx, arg)
	}

	return mockQueries.BaseMock.DeleteBook(ctx, arg)
}

func newTestUserID() uuid.UUID {
	return uuid.New()
}

func newTestBook(userId uuid.UUID) database.Book {
	return database.Book{
		ID:        uuid.New(),
		Title:     "The Great Go Debugger",
		Author:    "Gopher Max",
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestCreateBook(tTesting *testing.T) {
	userId := newTestUserID()
	testBook := newTestBook(userId)

	type createBookRequest struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	// 1. Success test case
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			CreateBookFunc: func(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
				if arg.UserID != userId {
					t.Errorf("Expected UserID %s, got %s", userId, arg.UserID)
				}

				return testBook, nil
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		requestBody, _ := json.Marshal(createBookRequest{
			Title:  testBook.Title,
			Author: testBook.Author,
		})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(requestBody))

		recorder := httptest.NewRecorder()

		apiConfig.CreateBook(recorder, request, userId)

		if recorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
		}
	})

	// 2. Invalid Input test case
	tTesting.Run("InvalidInput", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			CreateBookFunc: func(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
				t.Fatal("FATAL: CreateBook should not have been called on Invalid Input.")
				return database.Book{}, nil
			},
		}
		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		// Missing required field (Title)
		requestBody, _ := json.Marshal(createBookRequest{
			Author: "Invalid Author",
		})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		apiConfig.CreateBook(recorder, request, userId)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Database internal error test case
	tTesting.Run("InternalDBError", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			CreateBookFunc: func(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
				return database.Book{}, errors.New("simulated DB connection failure")
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		requestBody, _ := json.Marshal(createBookRequest{
			Title:  testBook.Title,
			Author: testBook.Author,
		})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		apiConfig.CreateBook(recorder, request, userId)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
		}
	})
}

func TestGetBook(tTesting *testing.T) {
	userId := newTestUserID()
	testBook := newTestBook(userId)

	// 1. Success test case
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return testBook, nil
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/books/%s", testBook.ID), nil)
		vars := map[string]string{
			"bookId": testBook.ID.String(),
		}
		request = mux.SetURLVars(request, vars) 
		recorder := httptest.NewRecorder()

		apiConfig.GetBook(recorder, request, userId)

		if recorder.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
		}

		var returnedBook struct{ ID uuid.UUID } // Simple struct for ID check

		if err := json.Unmarshal(recorder.Body.Bytes(), &returnedBook); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if returnedBook.ID != testBook.ID {
			t.Errorf("Expected book ID %s, got %s", testBook.ID, returnedBook.ID)
		}
	})

	// 2. Book Not Found test case
	tTesting.Run("BookNotFound", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return database.Book{}, sql.ErrNoRows
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		missingID := uuid.New().String()
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/books/%s", missingID), nil)
		vars := map[string]string{
			"bookId": missingID,
		}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.GetBook(recorder, request, userId)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Invalid Book ID Format
	tTesting.Run("InvalidIDFormat", func(t *testing.T) {
		mockQueries := &MockQueries{BaseMock: common.NewBaseMock()}
		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		request := httptest.NewRequest(http.MethodGet, "/api/v1/books/not-a-uuid", nil)
		vars := map[string]string{
			"bookId": "not-a-uuid",
		}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.GetBook(recorder, request, userId)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 4. Database Internal Error
	tTesting.Run("InternalDBError", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return database.Book{}, errors.New("simulated DB connection failure") // DB ERROR STUB
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/books/%s", testBook.ID), nil)
		vars := map[string]string{
			"bookId": testBook.ID.String(),
		}
		request = mux.SetURLVars(request, vars)

		apiConfig.GetBook(recorder, request, userId)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
		}
	})
}

func TestUpdateBook(tTesting *testing.T) {
	userId := newTestUserID()
	testBook := newTestBook(userId)

	type updateBookRequest struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	// Define the update payload
	updatePayload := updateBookRequest{
		Title:  "The Golang Manual (Updated)",
		Author: "Gopher Max (Revised)",
	}
	requestBody, _ := json.Marshal(updatePayload)
	
	updatedBook := testBook
	updatedBook.Title = updatePayload.Title
	updatedBook.Author = updatePayload.Author

	// 1. Success test case
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return testBook, nil
			},
			UpdateBookFunc: func(ctx context.Context, arg database.UpdateBookParams) (database.Book, error) {
				if arg.UserID != userId {
					t.Fatalf("Expected UserID %s in UpdateBook call, got %s", userId, arg.UserID)
				}
				if arg.Title != updatePayload.Title || arg.Author != updatePayload.Author {
					t.Fatal("UpdateBook was called with incorrect data")
				}
				return updatedBook, nil
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/books/%s", testBook.ID), bytes.NewBuffer(requestBody))
		
		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.UpdateBook(recorder, request, userId) 

		if recorder.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
		}
	})

	// 2. Book Not Found
	tTesting.Run("BookNotFound", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			UpdateBookFunc: func(ctx context.Context, arg database.UpdateBookParams) (database.Book, error) {
				return database.Book{}, sql.ErrNoRows 
			},
		}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		missingID := uuid.New()
		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/books/%s", missingID), bytes.NewBuffer(requestBody))
		
		vars := map[string]string{"bookId": missingID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.UpdateBook(recorder, request, userId)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Invalid input/body
	tTesting.Run("InvalidInput", func(t *testing.T) {
		mockQueries := &MockQueries{BaseMock: common.NewBaseMock()}

		apiConfig := BookAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		
		invalidBody, _ := json.Marshal(updateBookRequest{Author: "Only Author"})
		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/books/%s", testBook.ID), bytes.NewBuffer(invalidBody))
		
		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.UpdateBook(recorder, request, userId)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})
}
