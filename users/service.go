package users

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
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (userAPIConfig *UserAPIConfig) CreateUser(writer http.ResponseWriter, request *http.Request) {
	createUserParameters := CreateUserParameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&createUserParameters)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %s", decoderError))

		return
	}

	if strings.TrimSpace(createUserParameters.FirstName) == "" || strings.TrimSpace(createUserParameters.LastName) == "" || 
		strings.TrimSpace(createUserParameters.Email) == "" || strings.TrimSpace(createUserParameters.Password) == "" {
		common.ErrorResponse(writer, http.StatusBadRequest, "first_name, last_name, email and password fields are required")

		return
	}

	// Validate email.
	isEmailValid := common.IsEmailValid(createUserParameters.Email)

	if !isEmailValid {
		common.ErrorResponse(writer, http.StatusBadRequest, "error creating user: Invalid email address")

		return
	}

	// Check if email already exists.
	_, getUserError := userAPIConfig.DB.GetUserByEmail(request.Context(), createUserParameters.Email)

	if getUserError != nil {
		if getUserError != sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusInternalServerError, "failed to register. Please try again in a few minutes")

			return
		}
	} else {
		common.ErrorResponse(writer, http.StatusConflict, "failed to register. Email address already in use")

		return
	}

	// Validate password.
	isPasswordValid := common.IsPasswordValid(createUserParameters.Password)

	if !isPasswordValid {
		common.ErrorResponse(writer, http.StatusBadRequest, "error creating user: Invalid password. Password must contain at least 1 upper case letter, 1 lower case letter, 1 digit and must be 8 to 15 characters long.")

		return
	}

	hashedPassword, hashPasswordError := HashPassword(createUserParameters.Password)

	if hashPasswordError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("error creating user: %s", hashPasswordError))

		return
	}

	createUserParams := database.CreateUserParams {
		ID: uuid.New(),
		FirstName: createUserParameters.FirstName,
		LastName: createUserParameters.LastName,
		Email: createUserParameters.Email,
		Password: hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	newUser, createUserError := userAPIConfig.DB.CreateUser(request.Context(), createUserParams)

	if createUserError != nil {
		common.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("error creating user: %s", createUserError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseUserToUserJSON(newUser))
}

func (userAPIConfig *UserAPIConfig) Login(writer http.ResponseWriter, request *http.Request) {
	userLoginParameters := UserLoginParameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&userLoginParameters)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %s", decoderError))

		return
	}

	if strings.TrimSpace(userLoginParameters.Email) == "" || strings.TrimSpace(userLoginParameters.Password) == "" {
		common.ErrorResponse(writer, http.StatusBadRequest, "email and password fields are required")

		return
	}

	getUser, getUserError := userAPIConfig.DB.GetUserByEmail(request.Context(), userLoginParameters.Email)

	if getUserError != nil {
		if getUserError == sql.ErrNoRows {
			common.ErrorResponse(writer, http.StatusUnauthorized, "incorrect email address or password")
		} else {
			common.ErrorResponse(writer, http.StatusInternalServerError, "failed to login, Please try again in a few minutes")
		}	

		return
	}

	verifyPasswordError := VerifyPassword(userLoginParameters.Password, getUser.Password)

	if verifyPasswordError != nil {
		common.ErrorResponse(writer, http.StatusUnauthorized, "incorrect email address or password")

		return
	}

	// Private and public keys used the following settings for this project:
	// Curve: SECG secp256r1 / X9.62 prime256v1 / NIST P-256
	// Output Type: PEM text
	// Format: PKCS#8
	newToken := jwt.NewWithClaims(
		jwt.SigningMethodES256, 
		jwt.MapClaims{ 
			"email": getUser.Email, 
			"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	signedToken, signedStringError := newToken.SignedString(userAPIConfig.APIConfig.JWTSigningKey)
	
	if signedStringError != nil {
		log.Printf("signing error: %v", signedStringError)
		common.ErrorResponse(writer, http.StatusInternalServerError, "failed to login, Please try again in a few minutes")

		return
	}

	userAuthorized := DatabaseUserToUserAuthorizedJSON(getUser)
	userAuthorized.Token = signedToken

	common.JSONResponse(writer, http.StatusOK, userAuthorized)
}