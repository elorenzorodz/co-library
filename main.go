package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elorenzorodz/co-library/book_borrows"
	"github.com/elorenzorodz/co-library/books"
	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/elorenzorodz/co-library/user_subscribers"
	"github.com/elorenzorodz/co-library/users"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env.dev")

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

	// Users endpoints.
	userAPIConfig := users.UserAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(apiVersion + "/user/register", userAPIConfig.CreateUser).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/user/login", userAPIConfig.Login).Methods("POST")

	// Books endpoints.
	bookAPIConfig := books.BookAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(apiVersion + "/books", apiConfig.Authorization(bookAPIConfig.CreateBook)).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/books", apiConfig.Authorization(bookAPIConfig.GetBooks)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/books/browse", apiConfig.Authorization(bookAPIConfig.BrowseBooks)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/books/browse/{id}", apiConfig.Authorization(bookAPIConfig.BrowseBooksByUserID)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/books/{id}", apiConfig.Authorization(bookAPIConfig.GetBook)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/books/{id}", apiConfig.Authorization(bookAPIConfig.UpdateBook)).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/books/{id}", apiConfig.Authorization(bookAPIConfig.DeleteBook)).Methods("DELETE")

	// Book borrows endpoints.
	bookBorrowAPIConfig := book_borrows.BookBorrowAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(apiVersion + "/books/issue/{id}", apiConfig.Authorization(bookBorrowAPIConfig.IssueBook)).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/books/return/{id}", apiConfig.Authorization(bookBorrowAPIConfig.ReturnBook)).Methods("POST")

	// User subscrbers endpoints.
	userSubscriberAPIConfig := user_subscribers.UserSubscriberAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(apiVersion + "/users/subscribe/{user_id}", apiConfig.Authorization(userSubscriberAPIConfig.CreateUserSubscriber)).Methods("POST")
	muxRouter.HandleFunc(apiVersion + "/users/unsubscribe/{user_id}", apiConfig.Authorization(userSubscriberAPIConfig.DeleteUserSubscriber)).Methods("DELETE")
	muxRouter.HandleFunc(apiVersion + "/users/subscribers", apiConfig.Authorization(userSubscriberAPIConfig.GetUserSubscribers)).Methods("GET")
	muxRouter.HandleFunc(apiVersion + "/users/subscriptions", apiConfig.Authorization(userSubscriberAPIConfig.GetUserSubscriptions)).Methods("GET")

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