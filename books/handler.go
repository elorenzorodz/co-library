package books

import (
	"github.com/elorenzorodz/co-library/internal/database"
)

func DatabaseBookToBookJSON(databaseBook database.Book) Book {
	return Book{
		ID:        databaseBook.ID,
		Title:     databaseBook.Title,
		Author:    databaseBook.Author,
		CreatedAt: databaseBook.CreatedAt,
		UpdatedAt: databaseBook.UpdatedAt,
		UserID:    databaseBook.UserID,
	}
}

func DatabaseBooksToBooksJSON(databaseBooks []database.Book) []Book {
	books := []Book{}

	for _, databaseBook := range databaseBooks {
		books = append(books, DatabaseBookToBookJSON(databaseBook))
	}

	return books
}