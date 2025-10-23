package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elorenzorodz/co-library/book_borrows"
	"github.com/elorenzorodz/co-library/books"
	"github.com/elorenzorodz/co-library/common"
	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/elorenzorodz/co-library/middleware"
	"github.com/elorenzorodz/co-library/user_subscribers"
	"github.com/elorenzorodz/co-library/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if envFileLoadError := godotenv.Load(".env.dev"); envFileLoadError != nil {
		log.Fatal("Error loading .env file:", envFileLoadError)
	}

	envConfig := common.LoadEnvConfig()

	publicBytes, readPublicKeyError := os.ReadFile("public.pem")
	
    if readPublicKeyError != nil {
        log.Fatal("Error reading public.pem:", readPublicKeyError)
    }

    publicKey, parsingPublicKeyError := jwt.ParseECPublicKeyFromPEM(publicBytes)

    if parsingPublicKeyError != nil {
        log.Fatal("Error parsing public key:", parsingPublicKeyError)
    }

	privateBytes, readPrivateKeyError := os.ReadFile("private.pem")

	if readPrivateKeyError != nil {
		log.Fatal("private key read file error: ", readPrivateKeyError)
	}

	parsedPrivateKey, parsePrivateKeyError := jwt.ParseECPrivateKeyFromPEM(privateBytes)

	if parsePrivateKeyError != nil {
		log.Fatal("parse private key error: ", parsePrivateKeyError)
	}

	routeAPIPrefix := fmt.Sprintf("/api/%s", envConfig.APIVersion)

	dbConnection := common.OpenDBConnection(envConfig.DBUrl)

	database := database.New(dbConnection)

	apiConfig := common.APIConfig {
		DB: database,
		JWTValidationKey: publicKey,
		JWTSigningKey: parsedPrivateKey,
		MailgunAPIKey: envConfig.MailgunAPIKey,
		MailgunSendingDomain: envConfig.MailgunSendingDomain,
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
	muxRouter.HandleFunc(routeAPIPrefix + "/books/browse/{userId}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.BrowseBooksByUserID)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/{bookId}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.GetBook)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/{bookId}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.UpdateBook)).Methods("PATCH")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/{bookId}", middleware.Authorization(&bookAPIConfig.APIConfig, bookAPIConfig.DeleteBook)).Methods("DELETE")

	// Book borrows endpoints.
	bookBorrowAPIConfig := book_borrows.BookBorrowAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(routeAPIPrefix + "/books/issue/{bookId}", middleware.Authorization(&bookBorrowAPIConfig.APIConfig, bookBorrowAPIConfig.IssueBook)).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/books/return/{bookBorrowId}", middleware.Authorization(&bookBorrowAPIConfig.APIConfig, bookBorrowAPIConfig.ReturnBook)).Methods("POST")

	// User subscrbers endpoints.
	userSubscriberAPIConfig := user_subscribers.UserSubscriberAPIConfig {
		APIConfig: apiConfig,
	}

	muxRouter.HandleFunc(routeAPIPrefix + "/users/subscribe/{userId}", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.CreateUserSubscriber)).Methods("POST")
	muxRouter.HandleFunc(routeAPIPrefix + "/users/unsubscribe/{userId}", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.DeleteUserSubscriber)).Methods("DELETE")
	muxRouter.HandleFunc(routeAPIPrefix + "/users/subscribers", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.GetUserSubscribers)).Methods("GET")
	muxRouter.HandleFunc(routeAPIPrefix + "/users/subscriptions", middleware.Authorization(&userSubscriberAPIConfig.APIConfig, userSubscriberAPIConfig.GetUserSubscriptions)).Methods("GET")

	log.Printf("Server starting on port %v", envConfig.Port)

	server := &http.Server{
		Handler: muxRouter,
		Addr: ":" + envConfig.Port,
	}

	serverError := server.ListenAndServe()

	if serverError != nil {
		log.Fatal(serverError)
	}
}