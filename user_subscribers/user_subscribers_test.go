package user_subscribers

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

type MockUserSubscribersDB struct {
	*common.BaseMock

	GetUserByIDFunc func(ctx context.Context, id uuid.UUID) (database.User, error)

	GetUserSubscriberFunc    func(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error)
	CreateUserSubscriberFunc func(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error)
	DeleteUserSubscriberFunc func(ctx context.Context, arg database.DeleteUserSubscriberParams) (int64, error)
	GetUserSubscribersFunc   func(ctx context.Context, userID uuid.UUID) ([]database.UserSubscriber, error)
	GetUserSubscriptionsFunc func(ctx context.Context, subscriberID uuid.UUID) ([]database.UserSubscriber, error)
}

func (mockQueries *MockUserSubscribersDB) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	if mockQueries.GetUserByIDFunc != nil {
		return mockQueries.GetUserByIDFunc(ctx, id)
	}
	return mockQueries.BaseMock.GetUserByID(ctx, id)
}

func (mockQueries *MockUserSubscribersDB) GetUserSubscriber(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error) {
	if mockQueries.GetUserSubscriberFunc != nil {
		return mockQueries.GetUserSubscriberFunc(ctx, arg)
	}
	return mockQueries.BaseMock.GetUserSubscriber(ctx, arg)
}

func (mockQueries *MockUserSubscribersDB) CreateUserSubscriber(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error) {
	if mockQueries.CreateUserSubscriberFunc != nil {
		return mockQueries.CreateUserSubscriberFunc(ctx, arg)
	}
	return mockQueries.BaseMock.CreateUserSubscriber(ctx, arg) 
}

func (mockQueries *MockUserSubscribersDB) DeleteUserSubscriber(ctx context.Context, arg database.DeleteUserSubscriberParams) (int64, error) {
	if mockQueries.DeleteUserSubscriberFunc != nil {
		return mockQueries.DeleteUserSubscriberFunc(ctx, arg)
	}
	return mockQueries.BaseMock.DeleteUserSubscriber(ctx, arg) 
}

func (mockQueries *MockUserSubscribersDB) GetUserSubscribers(ctx context.Context, userID uuid.UUID) ([]database.UserSubscriber, error) {
	if mockQueries.GetUserSubscribersFunc != nil {
		return mockQueries.GetUserSubscribersFunc(ctx, userID)
	}
	return mockQueries.BaseMock.GetUserSubscribers(ctx, userID)
}

func (mockQueries *MockUserSubscribersDB) GetUserSubscriptions(ctx context.Context, subscriberID uuid.UUID) ([]database.UserSubscriber, error) {
	if mockQueries.GetUserSubscriptionsFunc != nil {
		return mockQueries.GetUserSubscriptionsFunc(ctx, subscriberID)
	}
	return mockQueries.BaseMock.GetUserSubscriptions(ctx, subscriberID)
}

func newTestUserID() uuid.UUID {
	return uuid.New()
}

func newTestUser(id uuid.UUID) database.User {
	return database.User{
		ID:        id,
		FirstName:  "John",
		LastName: "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newTestUserSubscriber(userID, subscriberID uuid.UUID) database.UserSubscriber {
	return database.UserSubscriber{
		ID:           uuid.New(),
		UserID:       userID,
		SubscriberID: subscriberID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestCreateUserSubscriber(tTesting *testing.T) {
	targetUserID := newTestUserID()
	subscriberID := newTestUserID()
	testSubscriber := newTestUserSubscriber(targetUserID, subscriberID)

	// 1. Success: User is successfully subscribed.
	tTesting.Run("Success", func(t *testing.T) {
		mockDB := &MockUserSubscribersDB{
			BaseMock: common.NewBaseMock(),
			GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (database.User, error) {
				return newTestUser(id), nil
			},
			GetUserSubscriberFunc: func(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error) {
				return database.UserSubscriber{}, sql.ErrNoRows
			},
			CreateUserSubscriberFunc: func(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error) {
				return testSubscriber, nil
			},
		}

		apiConfig := UserSubscriberAPIConfig{APIConfig: common.APIConfig{DB: mockDB}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/subscribers/users/%s", targetUserID), nil)
		vars := map[string]string{"userId": targetUserID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.CreateUserSubscriber(recorder, request, subscriberID)

		if recorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
		}
	})

	// 2. Failure: Invalid User ID Format
	tTesting.Run("InvalidUserIDFormat", func(t *testing.T) {
		mockDB := &MockUserSubscribersDB{BaseMock: common.NewBaseMock()}
		apiConfig := UserSubscriberAPIConfig{APIConfig: common.APIConfig{DB: mockDB}}

		request := httptest.NewRequest(http.MethodPost, "/api/v1/subscribers/users/not-a-uuid", nil)
		vars := map[string]string{"userId": "not-a-uuid"}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.CreateUserSubscriber(recorder, request, subscriberID)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Failure: Cannot Subscribe to Self
	tTesting.Run("SubscribeToSelf", func(t *testing.T) {
		mockDB := &MockUserSubscribersDB{BaseMock: common.NewBaseMock()}
		apiConfig := UserSubscriberAPIConfig{APIConfig: common.APIConfig{DB: mockDB}}

		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/subscribers/users/%s", subscriberID), nil)
		vars := map[string]string{"userId": subscriberID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.CreateUserSubscriber(recorder, request, subscriberID)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 4. Failure: Target User Not Found
	tTesting.Run("TargetUserNotFound", func(t *testing.T) {
		mockDB := &MockUserSubscribersDB{
			BaseMock: common.NewBaseMock(),
			GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (database.User, error) {
				return database.User{}, sql.ErrNoRows
			},
		}

		apiConfig := UserSubscriberAPIConfig{APIConfig: common.APIConfig{DB: mockDB}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/subscribers/users/%s", targetUserID), nil)
		vars := map[string]string{"userId": targetUserID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.CreateUserSubscriber(recorder, request, subscriberID)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
		}
	})

	// 5. Failure: Already Subscribed
	tTesting.Run("AlreadySubscribed", func(t *testing.T) {
		mockDB := &MockUserSubscribersDB{
			BaseMock: common.NewBaseMock(),
			GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (database.User, error) {
				return newTestUser(id), nil
			},
			GetUserSubscriberFunc: func(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error) {
				return testSubscriber, nil
			},
			CreateUserSubscriberFunc: func(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error) {
				t.Fatal("CreateUserSubscriber should not be called")
				return database.UserSubscriber{}, nil
			},
		}

		apiConfig := UserSubscriberAPIConfig{APIConfig: common.APIConfig{DB: mockDB}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/subscribers/users/%s", targetUserID), nil)
		vars := map[string]string{"userId": targetUserID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.CreateUserSubscriber(recorder, request, subscriberID)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 6. Failure: Internal Error on CreateUserSubscriber
	tTesting.Run("CreateDBError", func(t *testing.T) {
		mockDB := &MockUserSubscribersDB{
			BaseMock: common.NewBaseMock(),
			GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (database.User, error) {
				return newTestUser(id), nil
			},
			GetUserSubscriberFunc: func(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error) {
				return database.UserSubscriber{}, sql.ErrNoRows
			},
			CreateUserSubscriberFunc: func(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error) {
				return database.UserSubscriber{}, errors.New("simulated DB error on create")
			},
		}

		apiConfig := UserSubscriberAPIConfig{APIConfig: common.APIConfig{DB: mockDB}}
		request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/subscribers/users/%s", targetUserID), nil)
		vars := map[string]string{"userId": targetUserID.String()}
		request = mux.SetURLVars(request, vars)
		recorder := httptest.NewRecorder()

		apiConfig.CreateUserSubscriber(recorder, request, subscriberID)

		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
		}
	})
}
