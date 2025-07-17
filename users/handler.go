package users

import "github.com/elorenzorodz/co-library/internal/database"

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