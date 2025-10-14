package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elorenzorodz/co-library/book_borrows"
	"github.com/elorenzorodz/co-library/books"
	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/elorenzorodz/co-library/middleware"
	"github.com/elorenzorodz/co-library/user_subscribers"
	"github.com/elorenzorodz/co-library/users"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env.dev")

	apiVersion := common.GetEnvVariable("API_VERSION")
	routeAPIPrefix := fmt.Sprintf("/api/%s", apiVersion)

	dbConnectionString := common.GetDBConnectionSettings()
	dbConnection := common.OpenDBConnection(dbConnectionString)

	database := database.New(dbConnection)

	apiConfig := common.APIConfig {
		DB: database,
	}

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(routeAPIPrefix + "/ping", common.Pong).Methods("GET")

	// Users endpoints.
	userAPIConfig := users.UserAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(routeAPIPrefix + "/user/register", userAPIConfig.CreateUser).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/user/login", userAPIConfig.Login).Methods("POST")

	// Books endpoints.
	bookAPIConfig := books.BookAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(routeAPIPrefix + "/books", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.CreateBook)).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/books", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.GetBooks)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/browse", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.BrowseBooks)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/browse/{id}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.BrowseBooksByUserID)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/{id}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.GetBook)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/{id}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.UpdateBook)).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/{id}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.DeleteBook)).Methods("DELETE")

	// Book borrows endpoints.
	bookBorrowAPIConfig := book_borrows.BookBorrowAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(routeAPIPrefix + "/books/issue/{id}", middleware.Authorization(&bookBorrowAPIConfig.APIConfig, bookBorrowAPIConfig.IssueBook)).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/return/{id}", middleware.Authorization(&bookBorrowAPIConfig.APIConfig, bookBorrowAPIConfig.ReturnBook)).Methods("POST")

	// User subscrbers endpoints.
	userSubscriberAPIConfig := user_subscribers.UserSubscriberAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(routeAPIPrefix + "/users/subscribe/{user_id}", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.CreateUserSubscriber)).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/users/unsubscribe/{user_id}", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.DeleteUserSubscriber)).Methods("DELETE")
	muxRouter.HandleFunc(routeAPIPrefix + "/users/subscribers", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.GetUserSubscribers)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/users/subscriptions", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.GetUserSubscriptions)).Methods("GET")

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