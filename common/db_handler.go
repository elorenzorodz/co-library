package common

import (
	"database/sql"
	"log"
)

func OpenDBConnection(dbUrl string) *sql.DB {
	connection, connectionError := sql.Open("postgres", dbUrl)

	if connectionError != nil {
		log.Fatal("Can't open connection to database:", connectionError)
	}

	pingError := connection.Ping()

	if pingError != nil {
		log.Fatal("Can't connect to database:", pingError)
	}

	return connection
}