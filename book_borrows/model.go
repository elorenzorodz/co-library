package book_borrows

import (
	"database/sql"
	"time"

	"github.com/elorenzorodz/co-library/common"
	"github.com/google/uuid"
)

type BookBorrowAPIConfig struct {
	common.APIConfig
}

type BookBorrow struct {
	ID uuid.UUID `json:"id"`
	IssuedAt time.Time `json:"issuedAt"`
	ReturnedAt sql.NullTime `json:"returnedAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	BookID uuid.UUID `json:"book_id"`
	BorrowerID uuid.UUID `json:"borrower_id"`
}
