package common

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
)

func IsEmailValid(email string) bool {
	emailRegex, emailValidationError := regexp.MatchString(`^[A-Za-z0-9]+([._\-][A-Za-z0-9]+)*@[A-Za-z0-9]+([\-\.][A-Za-z0-9]+)*\.[A-Za-z]{2,15}$`, email)

	if emailValidationError != nil {
		log.Printf("Invalid email: %s", emailValidationError)
	}

	return emailRegex
}

func IsPasswordValid(password string) bool {
	if len(password) < 8 || len(password) > 15 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpace bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsSpace(char):
			hasSpace = true
		}
	}

	return hasUpper && hasLower && hasDigit && !hasSpace
}

func ValidateJWTAndGetEmailClaim(signedToken string, publicKey interface{}) (string, error) {
	parsedToken, parsedTokenError := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return publicKey, nil
	})

	if parsedTokenError != nil {
		log.Printf("Token parse error: %s", parsedTokenError)

		return "", fmt.Errorf("token parse error: %s", parsedTokenError)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		emailClaim, exists := claims["email"]

		if !exists {
			log.Println("Email claim not found")

			return "", errors.New("invalid token")
		}

		email, ok := emailClaim.(string)

		if !ok {
			log.Println("Email claim is not a string")

			return "", errors.New("invalid token")
		}

		return email, nil
	}

	return "", errors.New("invalid token")
}