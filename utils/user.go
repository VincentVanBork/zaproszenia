package utils

import (
	"cloud.google.com/go/datastore"
	"context"
	"log"
	"zaproszenia/models"
)

func GetUsers(client *datastore.Client) map[string]string {
	log.Println("Getting Users")
	ctx := context.Background()
	query := datastore.NewQuery("User")
	var users []*models.User

	if _, err := client.GetAll(ctx, query, &users); err != nil {
		log.Println("Error getting all users", err)
		return map[string]string{}
	}

	if len(users) == 0 {
		log.Println("There are no users in users list")
		return map[string]string{}
	}
	return map[string]string{users[0].Username: users[0].Password}
}
