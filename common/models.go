package common

import "github.com/elorenzorodz/co-library/internal/database"

type EnvConfig struct {
	APIVersion string
	Port       string
	DBUrl      string
}

type APIConfig struct {
	DB        *database.Queries
	PublicKey interface{}
}