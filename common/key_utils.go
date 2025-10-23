package common

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func LoadAuthKeys() (signingKey interface{}, validationKey interface{}) {
	publicBytes, readPublicKeyError := os.ReadFile("public.pem")

	if readPublicKeyError != nil {
		log.Fatal("error reading public.pem:", readPublicKeyError)
	}

	parsedPublicKey, parsingPublicKeyError := jwt.ParseECPublicKeyFromPEM(publicBytes)

	if parsingPublicKeyError != nil {
		log.Fatal("error parsing public key:", parsingPublicKeyError)
	}

	privateBytes, readPrivateKeyError := os.ReadFile("private.pem")

	if readPrivateKeyError != nil {
		log.Fatal("private key read file error: ", readPrivateKeyError)
	}

	parsedPrivateKey, parsePrivateKeyError := jwt.ParseECPrivateKeyFromPEM(privateBytes)

	if parsePrivateKeyError != nil {
		log.Fatal("parse private key error: ", parsePrivateKeyError)
	}

	return parsedPublicKey, parsedPrivateKey
}