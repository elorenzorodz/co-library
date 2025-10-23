package common

import "github.com/elorenzorodz/co-library/internal/database"

type EnvConfig struct {
	APIVersion           string
	Port                 string
	DBUrl                string
	MailgunAPIKey        string
	MailgunSendingDomain string
}

type APIConfig struct {
	DB                   *database.Queries
	JWTValidationKey	 interface{}
	JWTSigningKey        interface{}
	MailgunAPIKey        string
	MailgunSendingDomain string
}