package common

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func Pong(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Message string
	}

	params := parameters{
		Message: "pong",
	}

	JSONResponse(writer, http.StatusOK, params)
}

func GetJWT(headers http.Header) (string, error){
	authorizationHeader := headers.Get("Authorization")

	if authorizationHeader == "" {
		return "", errors.New("no authentication info found")
	}

	authorizationHeaderValues := strings.Split(authorizationHeader, " ")

	if len(authorizationHeaderValues) != 2 {
		return "", errors.New("malformed authentication header")
	}

	if authorizationHeaderValues[0] != "Bearer" {
		return "", errors.New("malformed first part of authentication header")
	}

	return authorizationHeaderValues[1], nil
}

func JSONResponse(writer http.ResponseWriter, code int, payload interface{}) {
	data, jsonMarshalError := json.Marshal(payload)

	if jsonMarshalError != nil {
		log.Printf("failed to marshal JSON response: %v", payload)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(data)
}

func ErrorResponse(writer http.ResponseWriter, code int, message string) {
	// Log to server on 5XX status codes.
	if code > 499 {
		log.Println("responding with 5XX error:", message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	errResponse := errorResponse{
		Error: message,
	}

	JSONResponse(writer, code, errResponse)
}