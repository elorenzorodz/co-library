package common

import (
	"context"
	"database/sql"

	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/google/uuid"
)

type UserMock struct{}

func (m *UserMock) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	panic("CreateUser not implemented for this test (BaseMock)")
}

func (m *UserMock) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return database.User{}, sql.ErrNoRows
}

func (m *UserMock) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return database.User{}, sql.ErrNoRows
}

type BookMock struct{}

func (m *BookMock) CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
	panic("CreateBook not implemented for this test (BaseMock)")
}

func (m *BookMock) GetBook(ctx context.Context, id uuid.UUID) (database.Book, error) {
	return database.Book{}, sql.ErrNoRows
}

func (m *BookMock) GetBooks(ctx context.Context, userID uuid.UUID) ([]database.Book, error) {
	return []database.Book{}, nil
}

func (m *BookMock) BrowseBooks(ctx context.Context) ([]database.Book, error) {
	return []database.Book{}, nil
}

func (m *BookMock) UpdateBook(ctx context.Context, arg database.UpdateBookParams) (database.Book, error) {
	panic("UpdateBook not implemented for this test (BaseMock)")
}

func (m *BookMock) DeleteBook(ctx context.Context, arg database.DeleteBookParams) (int64, error) {
	return 0, nil
}

type BookBorrowMock struct{}

func (m *BookBorrowMock) GetBookBorrow(ctx context.Context, bookID uuid.UUID) (database.BookBorrow, error) {
	return database.BookBorrow{}, sql.ErrNoRows
}

func (m *BookBorrowMock) IssueBook(ctx context.Context, arg database.IssueBookParams) (database.BookBorrow, error) {
	panic("IssueBook not implemented for this test (BaseMock)")
}

func (m *BookBorrowMock) ReturnBook(ctx context.Context, arg database.ReturnBookParams) (database.BookBorrow, error) {
	panic("ReturnBook not implemented for this test (BaseMock)")
}

type UserSubscriberMock struct{}

func (m *UserSubscriberMock) CreateUserSubscriber(ctx context.Context, arg database.CreateUserSubscriberParams) (database.UserSubscriber, error) {
	panic("CreateUserSubscriber not implemented for this test (BaseMock)")
}

func (m *UserSubscriberMock) GetUserSubscriber(ctx context.Context, arg database.GetUserSubscriberParams) (database.UserSubscriber, error) {
	return database.UserSubscriber{}, sql.ErrNoRows
}

func (m *UserSubscriberMock) GetUserSubscribers(ctx context.Context, userID uuid.UUID) ([]database.UserSubscriber, error) {
	return []database.UserSubscriber{}, nil
}

func (m *UserSubscriberMock) GetUserSubscriptions(ctx context.Context, subscriberID uuid.UUID) ([]database.UserSubscriber, error) {
	return []database.UserSubscriber{}, nil
}

func (m *UserSubscriberMock) GetUsersBySubscriberID(ctx context.Context, userID uuid.UUID) ([]database.User, error) {
	return []database.User{}, nil
}

func (m *UserSubscriberMock) DeleteUserSubscriber(ctx context.Context, arg database.DeleteUserSubscriberParams) (int64, error) {
	return 0, nil
}

type BaseMock struct {
	*UserMock
	*BookMock
	*BookBorrowMock
	*UserSubscriberMock
}

func NewBaseMock() *BaseMock {
	return &BaseMock{
		UserMock:           &UserMock{},
		BookMock:           &BookMock{},
		BookBorrowMock:     &BookBorrowMock{},
		UserSubscriberMock: &UserSubscriberMock{},
	}
}