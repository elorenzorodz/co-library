package user_subscribers

import (
	"github.com/elorenzorodz/co-library/internal/database"
)

func DatabaseUserSubscriberToUserSubscriberJSON(databaseUserSubscriber database.UserSubscriber) UserSubscriber {
	return UserSubscriber{
		ID:            databaseUserSubscriber.ID,
		CreatedAt:     databaseUserSubscriber.CreatedAt,
		UpdatedAt:     databaseUserSubscriber.UpdatedAt,
		UserID:        databaseUserSubscriber.UserID,
		SubscriberID: databaseUserSubscriber.SubscriberID,
	}
}

func DatabaseUserSubscribersToUserSubscribersJSON(databaseUserSubscribers []database.UserSubscriber) []UserSubscriber {
	userSubscribers := []UserSubscriber{}

	for _, databaseUserSubscriber := range databaseUserSubscribers {
		userSubscribers = append(userSubscribers, DatabaseUserSubscriberToUserSubscriberJSON(databaseUserSubscriber))
	}

	return userSubscribers
}