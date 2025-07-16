package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elorenzorodz/co-library/common"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	apiVersion := common.GetEnvVariable("API_VERSION")
	apiVersion = fmt.Sprintf("/%s", apiVersion)

	// dbConnectionString := common.GetDBConnectionSettings()
	// dbConnection := openDBConnection(dbConnectionString)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(apiVersion + "/ping", common.Pong).Methods("GET")

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