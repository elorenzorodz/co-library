package users

import (
	"github.com/elorenzorodz/co-library/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func DatabaseUserToUserJSON(databaseUser database.User) User {
	return User{
		ID:        databaseUser.ID,
		FirstName: databaseUser.FirstName,
		LastName: databaseUser.LastName,
		Email: databaseUser.Email,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
	}
}

func DatabaseUserToUserAuthorizedJSON(databaseUser database.User) UserAuthorized {
	return UserAuthorized{
		Email: databaseUser.Email,
	}
}

func HashPassword(password string) (string, error) {
	bytes, hashPasswordError := bcrypt.GenerateFromPassword([]byte(password), 14)

    return string(bytes), hashPasswordError
}

func VerifyPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}