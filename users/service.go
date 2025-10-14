package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", decoderError))

		return
	}

	// Validate email.
	isEmailValid := common.IsEmailValid(createUserParameters.Email)

	if !isEmailValid {
		common.ErrorResponse(writer, http.StatusBadRequest, "Error creating user: Invalid email address")

		return
	}

	// Validate password.
	isPasswordValid := common.IsPasswordValid(createUserParameters.Password)

	if !isPasswordValid {
		common.ErrorResponse(writer, http.StatusBadRequest, "Error creating user: Invalid password. Password must contain at least 1 upper case letter, 1 lower case letter, 1 digit and must be 8 to 15 characters long.")

		return
	}

	hashedPassword, hashPasswordError := HashPassword(createUserParameters.Password)

	if hashPasswordError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error creating user: %s", hashPasswordError))

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
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error creating user: %s", createUserError))

		return
	}

	common.JSONResponse(writer, http.StatusCreated, DatabaseUserToUserJSON(newUser))
}

func (userAPIConfig *UserAPIConfig) Login(writer http.ResponseWriter, request *http.Request) {
	userLoginParameters := UserLoginParameters{}

	decoder := json.NewDecoder(request.Body)
	decoderError := decoder.Decode(&userLoginParameters)

	if decoderError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", decoderError))

		return
	}

	getUser, getUserError := userAPIConfig.DB.GetUserByEmail(request.Context(), userLoginParameters.Email)

	if getUserError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", getUserError))

		return
	}

	verifyPasswordError := VerifyPassword(userLoginParameters.Password, getUser.Password)

	if verifyPasswordError != nil {
		common.ErrorResponse(writer, http.StatusBadRequest, "Incorrect email address or password")

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

	bytes, readFileError := os.ReadFile("private.pem")

	if readFileError != nil {
		log.Printf("Read file error %v", readFileError)
		common.ErrorResponse(writer, http.StatusBadRequest, "Failed to login. Please try again in a few minutes")

		return
	}

	key, parsePrivateKeyError := jwt.ParseECPrivateKeyFromPEM(bytes)

	if parsePrivateKeyError != nil {
		log.Printf("Parse error %v", parsePrivateKeyError)
		common.ErrorResponse(writer, http.StatusBadRequest, "Failed to login. Please try again in a few minutes")

		return
	}

	signedToken, signedStringError := newToken.SignedString(key)

	if signedStringError != nil {
		log.Printf("Signing error %v", signedStringError)
		common.ErrorResponse(writer, http.StatusBadRequest, "Failed to login. Please try again in a few minutes")

		return
	}

	userAuthorized := DatabaseUserToUserAuthorizedJSON(getUser)
	userAuthorized.Token = signedToken

	common.JSONResponse(writer, http.StatusOK, userAuthorized)
}