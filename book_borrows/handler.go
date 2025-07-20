package book_borrows

import (
	"github.com/elorenzorodz/co-library/internal/database"
)

func DatabaseBookBorrowToBookBorrowJSON(databaseBookBorrow database.BookBorrow) BookBorrow {
	return BookBorrow{
		ID:        databaseBookBorrow.ID,
		IssuedAt:     databaseBookBorrow.IssuedAt,
		ReturnedAt:    databaseBookBorrow.ReturnedAt,
		CreatedAt: databaseBookBorrow.CreatedAt,
		UpdatedAt: databaseBookBorrow.UpdatedAt,
		BookID:    databaseBookBorrow.BookID,
		BorrowerID:    databaseBookBorrow.BorrowerID,
	}
}