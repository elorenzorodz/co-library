package common

import (
	"context"

	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
)

type EnvConfig struct {
	APIVersion           string
	Port                 string
	DBUrl                string
	MailgunAPIKey        string
	MailgunSendingDomain string
}

type APIConfig struct {
	DB                   Querier
	JWTValidationKey     interface{}
	JWTSigningKey        interface{}
	MailgunAPIKey        string
	MailgunSendingDomain string
}

type Querier interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error)

	CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error)
	GetBook(ctx context.Context, id uuid.UUID) (database.Book, error)
	GetBooks(ctx context.Context, userID uuid.UUID) ([]database.Book, error)
	BrowseBooks(ctx context.Context) ([]database.Book, error)
	UpdateBook(ctx context.Context, arg database.UpdateBookParams) (database.Book, error)
	DeleteBook(ctx context.Context, arg database.DeleteBookParams) (int64, error)

	GetBookBorrow(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error)
	IssueBook(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error)
	ReturnBook(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error)

	CreateUserSubscriber(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error)
	GetUserSubscriber(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error)
	GetUserSubscribers(ctx context.Context, userID uuid.UUID) ([]database.UserSubscriber, error)
	GetUserSubscriptions(ctx context.Context, subscriberID uuid.UUID) ([]database.UserSubscriber, error)
	GetUsersBySubscriberID(ctx context.Context, userID uuid.UUID) ([]database.User, error)
	DeleteUserSubscriber(ctx context.Context, arg database.DeleteUserSubscriberParams) (int64, error)
}