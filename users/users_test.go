package users

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)



type MockQueries struct {
    *common.BaseMock

    CreateUserFunc     func(ctx context.Context, arg database.CreateUserParams) (database.User, error)
    GetUserByEmailFunc func(ctx context.Context, email string) (database.User, error)
    GetUserByIDFunc    func(ctx context.Context, id uuid.UUID) (database.User, error)
}

func (mockQueries *MockQueries) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	if mockQueries.CreateUserFunc == nil {
        return mockQueries.BaseMock.CreateUser(ctx, arg)
	}
	return mockQueries.CreateUserFunc(ctx, arg)
}

func (mockQueries *MockQueries) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	if mockQueries.GetUserByEmailFunc != nil {
		return mockQueries.GetUserByEmailFunc(ctx, email)
	}
    
	return mockQueries.BaseMock.GetUserByEmail(ctx, email)
}

func (mockQueries *MockQueries) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	if mockQueries.GetUserByIDFunc != nil {
		return mockQueries.GetUserByIDFunc(ctx, id)
	}
	
	return mockQueries.BaseMock.GetUserByID(ctx, id)
}

func newTestUser() database.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("!Password123"), bcrypt.DefaultCost)
	return database.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     "test@email.com",
		Password:  string(hashedPassword),
	}
}

func TestCreateUser(tTesting *testing.T) {
	testUser := newTestUser()

	// Success test case.
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			CreateUserFunc: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
				return testUser, nil
			},
			GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
				return database.User{}, sql.ErrNoRows
			},
		}

		userAPIConfig := UserAPIConfig{
			APIConfig: common.APIConfig{DB: mockQueries},
		}

		requestBody, _ := json.Marshal(struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			Password  string `json:"password"`
		}{FirstName: testUser.FirstName, LastName: testUser.LastName, Email: testUser.Email, Password: "!Password123"})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		userAPIConfig.CreateUser(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
		}
	})

	// 2. INVALID INPUT Test Case
	tTesting.Run("InvalidInput", func(t *testing.T) {
		mockQueries := &MockQueries{ BaseMock: common.NewBaseMock() }
		userAPIConfig := UserAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		// Missing required field (password)
		requestBody, _ := json.Marshal(struct {
			Email string `json:"email"`
		}{Email: "invalid@example.com"})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		userAPIConfig.CreateUser(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Duplicate email test case (simulated DB constraint error)
	tTesting.Run("DuplicateEmail", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
				return testUser, nil
			},
			CreateUserFunc: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
				t.Fatal("FATAL: Handler ignored existing user and called CreateUser.")
				return database.User{}, nil
			},
		}

		userAPIConfig := UserAPIConfig{APIConfig: common.APIConfig{DB: mockQueries}}

		requestBody, _ := json.Marshal(struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			Password  string `json:"password"`
		}{FirstName: testUser.FirstName, LastName: testUser.LastName, Email: testUser.Email, Password: "!Password123"})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		userAPIConfig.CreateUser(recorder, request)

		if recorder.Code != http.StatusConflict {
			t.Errorf("Expected status %d (Conflict), got %d. Body: %s", http.StatusConflict, recorder.Code, recorder.Body.String())
		}
	})
}

func TestLogin(tTesting *testing.T) {
	validPassword := "!Password123"
	// Hash the correct password for the mock user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	testUser := database.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     "login@example.com",
		Password:  string(hashedPassword),
	}

	// This is a test private key only.
	privateKey, _ := jwt.ParseECPrivateKeyFromPEM([]byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg7zcfAR2TkjaT+h2M
MoI7ZWp3iJEgZSM8l50WrXaoFKuhRANCAAQDHEeHaLjjx7k5XQ23iTnGWSgwWIuK
GEPIHTMy6cJHSf+xLGKcIp40vHg1A9Rg8GeWhax4bIghE3cuKj5RNyQc
-----END PRIVATE KEY-----
`))

	// 1. Success test case
	tTesting.Run("Success", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
				return testUser, nil
			},
		}

		userAPIConfig := UserAPIConfig{APIConfig: common.APIConfig{DB: mockQueries, JWTSigningKey: privateKey}}

		requestBody, _ := json.Marshal(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{Email: testUser.Email, Password: validPassword})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		userAPIConfig.Login(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
		}

		// Optional: Verify JWT token structure in response body
		var resp struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Could not unmarshal response body: %v", err)
		}
		if resp.Token == "" {
			t.Error("Expected a token in the response")
		}
	})

	// 2. User not found test case
	tTesting.Run("UserNotFound", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
				return database.User{}, sql.ErrNoRows
			},
		}

		userAPIConfig := UserAPIConfig{APIConfig: common.APIConfig{DB: mockQueries, JWTSigningKey: privateKey}}

		requestBody, _ := json.Marshal(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{Email: "nonexistent@example.com", Password: validPassword})

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		userAPIConfig.Login(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
		}
	})

	// 3. Invalid credentials test case (wrong password)
	tTesting.Run("InvalidCredentials", func(t *testing.T) {
		mockQueries := &MockQueries{
			BaseMock: common.NewBaseMock(),
			GetUserByEmailFunc: func(ctx context.Context, email string) (database.User, error) {
				return testUser, nil // User is found
			},
		}

		userAPIConfig := UserAPIConfig{APIConfig: common.APIConfig{DB: mockQueries, JWTSigningKey: privateKey}}

		requestBody, _ := json.Marshal(struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{Email: testUser.Email, Password: "WrongPassword456"}) // Wrong password

		request := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(requestBody))
		recorder := httptest.NewRecorder()

		userAPIConfig.Login(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, recorder.Code, recorder.Body.String())
		}
	})
}
