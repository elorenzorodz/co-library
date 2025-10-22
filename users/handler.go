package users

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/elorenzorodz/co-library/internal/database"
	"github.com/mailgun/mailgun-go/v4"
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

func DispatchNewBookAlertsSync(bookTitle string, subscribers []database.User, sender database.User, mailGunAPIKey string, mailGunSendingDomain string) {
	waitGroup := &sync.WaitGroup{}

	for _, subscriber := range subscribers {
		waitGroup.Add(1)

		senderName := fmt.Sprintf("%s %s", sender.FirstName, sender.LastName)
		subscriberName := fmt.Sprintf("%s %s", subscriber.FirstName, subscriber.LastName)

		go SendNewBookAlert(mailGunSendingDomain, mailGunAPIKey, senderName, sender.Email, subscriberName, subscriber.Email, bookTitle, waitGroup)
	}

	waitGroup.Wait()

	log.Printf("New book alert sent to %v subscribers", len(subscribers))
}

func SendNewBookAlert(mailGunSendingDomain, mailGunAPIKey, senderName, senderEmail, subscriberName, subscriberEmail, bookTitle string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	fromNameAndEmail := fmt.Sprintf("%s <%s>", senderName, senderEmail)
	toNameAndEmail := fmt.Sprintf("%s <%s>", subscriberName, subscriberEmail)

	mg := mailgun.NewMailgun(mailGunSendingDomain, mailGunAPIKey)

	mailgunMessage := mailgun.NewMessage(
		fromNameAndEmail,
		"My Library Just Got Updated",
		fmt.Sprintf("Hi %s, \n\nI've added a new book in my library: %s \n\nCheck it out! Thank you.", subscriberName, bookTitle),
		toNameAndEmail,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	sendMessage, id, sendError := mg.Send(ctx, mailgunMessage)
	
	if sendError != nil {
		log.Printf("Mailgun send error | ID: %s | Message: %s | Error: %s", id, sendMessage, sendError)
	}
}