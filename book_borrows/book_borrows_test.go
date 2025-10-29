package book_borrows

import (
	"context"
	"database/sql"
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

	GetBookFunc       func(ctx context.Context, id uuid.UUID) (database.Book, error)
	GetBookBorrowFunc func(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error)
	IssueBookFunc     func(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error)
	ReturnBookFunc    func(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error)
}

func (mockQueries *MockQueries) GetBook(ctx context.Context, id uuid.UUID) (database.Book, error) {
	if mockQueries.GetBookFunc != nil {
		return mockQueries.GetBookFunc(ctx, id)
	}
	
	return mockQueries.BaseMock.GetBook(ctx, id)
}

func (mockQueries *MockQueries) GetBookBorrow(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error) {
	if mockQueries.GetBookBorrowFunc != nil {
		return mockQueries.GetBookBorrowFunc(ctx, bookID)
	}
	
	return mockQueries.BaseMock.GetBookBorrow(ctx, bookID)
}

func (mockQueries *MockQueries) IssueBook(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error) {
	if mockQueries.IssueBookFunc != nil {
		return mockQueries.IssueBookFunc(ctx, arg)
	}
	
	return mockQueries.BaseMock.IssueBook(ctx, arg)
}

func (mockQueries *MockQueries) ReturnBook(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error) {
	if mockQueries.ReturnBookFunc != nil {
		return mockQueries.ReturnBookFunc(ctx, arg)
	}
	
	return mockQueries.BaseMock.ReturnBook(ctx, arg)
}

func newTestUserID() uuid.UUID {
	return uuid.New()
}

func newTestBook(userId uuid.UUID) database.Book {
	return database.Book{
		ID:        uuid.New(),
		Title:     "The Go Test Manual",
		Author:    "Test Author",
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newTestBookBorrow(bookID, borrowerID uuid.UUID) database.BookBorrow {
	return database.BookBorrow{
		ID:         uuid.New(),
		BookID:     bookID,
		BorrowerID: borrowerID,
		IssuedAt:   time.Now(),
		ReturnedAt: sql.NullTime{Valid: false},
	}
}

func TestIssueBook(tTesting *testing.T) {
	bookUserId := newTestUserID()
	borrowerID := newTestUserID()
	testBook := newTestBook(bookUserId)
	testBorrow := newTestBookBorrow(testBook.ID, borrowerID)

	base := common.NewBaseMock()

	// 1. Success: book is available and successfully issued.
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return testBook, nil
			},
			GetBookBorrowFunc: func(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error) {
				return database.BookBorrow{}, sql.ErrNoRows
			},
			IssueBookFunc: func(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error) {
				if arg.BookID != testBook.ID || arg.BorrowerID != borrowerID {
					t.Fatalf("IssueBook called with wrong IDs")
				}
				return testBorrow, nil
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/book_borrows/issue/%s", testBook.ID), nil)

		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.IssueBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
		}
	})

	// 2. Failure: invalid book ID format
	tTesting.Run("InvalidBookIDFormat", func(t *testing.T) {
		mockQueries := &MockQueries{BaseMock: base}
		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/book_borrows/issue/not-a-uuid", nil)
		vars := map[string]string{"bookId": "not-a-uuid"}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.IssueBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Failure: book not found
	tTesting.Run("BookNotFound", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return database.Book{}, sql.ErrNoRows
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/book_borrows/issue/%s", testBook.ID), nil)

		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.IssueBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
		}
	})

	// 4. Failure: internal error on GetBook
	tTesting.Run("GetBookDBError", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return database.Book{}, errors.New("simulated DB error on get book") // DB error
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/book_borrows/issue/%s", testBook.ID), nil)

		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.IssueBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
		}
	})

	// 5. Failure: borrower is the book owner
	tTesting.Run("BorrowerIsOwner", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return testBook, nil
			},
			IssueBookFunc: func(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error) {
				t.Fatal("IssueBook should not be called when borrower is owner")
				return database.BookBorrow{}, nil
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/book_borrows/issue/%s", testBook.ID), nil)

		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		// Use the book owner's ID as the borrower ID.
		apiConfig.IssueBook(recorder, request, bookUserId)

		if recorder.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusForbidden, recorder.Code, recorder.Body.String())
		}
	})

	// 6. Failure: book already issued
	tTesting.Run("BookAlreadyIssued", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return testBook, nil
			},
			GetBookBorrowFunc: func(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error) {
				return newTestBookBorrow(bookID, newTestUserID()), nil
			},
			IssueBookFunc: func(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error) {
				t.Fatal("IssueBook should not be called when book is already issued")
				return database.BookBorrow{}, nil
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/book_borrows/issue/%s", testBook.ID), nil)

		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.IssueBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, recorder.Code, recorder.Body.String())
		}
	})

	// 7. Failure: internal error on IssueBook creation
	tTesting.Run("IssueBookDBError", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			GetBookFunc: func(ctx context.Context, id uuid.UUID) (database.Book, error) {
				return testBook, nil
			},
			GetBookBorrowFunc: func(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error) {
				return database.BookBorrow{}, sql.ErrNoRows
			},
			IssueBookFunc: func(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error) {
				return database.BookBorrow{}, errors.New("simulated DB error on issue book")
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/book_borrows/issue/%s", testBook.ID), nil)

		vars := map[string]string{"bookId": testBook.ID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.IssueBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
		}
	})
}

func TestReturnBook(tTesting *testing.T) {
	borrowerID := newTestUserID()
	bookID := uuid.New()
	testBorrowID := uuid.New()
	returnedBorrow := newTestBookBorrow(bookID, borrowerID)
	returnedBorrow.ReturnedAt.Valid = true
	returnedBorrow.ReturnedAt.Time = time.Now()

	base := common.NewBaseMock()

	// 1. Success: book is successfully returned.
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			ReturnBookFunc: func(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error) {
				if arg.ID != testBorrowID || arg.BorrowerID != borrowerID {
					t.Fatalf("ReturnBook called with wrong IDs")
				}
				return returnedBorrow, nil
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/book_borrows/return/%s", testBorrowID), nil)

		vars := map[string]string{"bookBorrowId": testBorrowID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.ReturnBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
		}
	})

	// 2. Failure: invalid borrow ID format
	tTesting.Run("InvalidBorrowIDFormat", func(t *testing.T) {
		mockQueries := &MockQueries{BaseMock: base}
		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		request := httptest.NewRequest(http.MethodPatch, "/api/v1/book_borrows/return/not-a-uuid", nil)
		vars := map[string]string{"bookBorrowId": "not-a-uuid"}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.ReturnBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Failure: not found, unauthorized, or already returned (sql.ErrNoRows)
	tTesting.Run("NotFoundOrUnauthorized", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			ReturnBookFunc: func(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error) {
				return database.BookBorrow{}, sql.ErrNoRows
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/book_borrows/return/%s", testBorrowID), nil)

		vars := map[string]string{"bookBorrowId": testBorrowID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		// Note: The handler returns 400 Bad Request for this specific DB error
		apiConfig.ReturnBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 4. Failure: Internal Error on ReturnBook
	tTesting.Run("ReturnBookDBError", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: base,
			ReturnBookFunc: func(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error) {
				return database.BookBorrow{}, errors.New("simulated DB error on return book") // DB error
			},
		}

		apiConfig := BookBorrowAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}
		request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/book_borrows/return/%s", testBorrowID), nil)

		vars := map[string]string{"bookBorrowId": testBorrowID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.ReturnBook(recorder, request, borrowerID)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
		}
	})
}