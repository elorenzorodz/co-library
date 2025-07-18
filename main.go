package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elorenzorodz/co-library/books"
	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/elorenzorodz/co-library/users"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	apiVersion := common.GetEnvVariable("API_VERSION")
	apiVersion = fmt.Sprintf("/%s", apiVersion)

	dbConnectionString := common.GetDBConnectionSettings()
	dbConnection := common.OpenDBConnection(dbConnectionString)

	database := database.New(dbConnection)

	apiConfig := common.APIConfig {
		DB: database,
	}

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(apiVersion + "/ping", common.Pong).Methods("GET")

	userAPIConfig := users.UserAPIConfig {
		APIConfig: apiConfig,
	}

	// Users endpoints.
	muxRouter.HandleFunc(apiVersion + "/user", userAPIConfig.CreateUser).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/user/login", userAPIConfig.Login).Methods("POST")

	bookAPIConfig := books.BookAPIConfig {
		APIConfig: apiConfig,
	}

	// Books endpoints.
	muxRouter.HandleFunc(apiVersion + "/books", bookAPIConfig.Authorization(bookAPIConfig.CreateBook)).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/books", bookAPIConfig.Authorization(bookAPIConfig.GetBooks)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/books/{id}", bookAPIConfig.Authorization(bookAPIConfig.GetBook)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/books/{id}", bookAPIConfig.Authorization(bookAPIConfig.UpdateBook)).Methods("POST")

	http.Handle("/", muxRouter)

	port := common.GetEnvVariable("PORT")

	log.Printf("Server starting on port %v", port)

	server := &http.Server{
		Handler: nil,
		Addr: ":" + port,
	}

	serverError := server.ListenAndServe()

	if serverError != nil {
		log.Fatal(serverError)
	}
}