// Sample datastore demonstrates use of the cloud.google.com/go/datastore package from App Engine flexible.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"zaproszenia/models"

	"cloud.google.com/go/datastore"
)

var datastoreClient *datastore.Client

func main() {
	ctx := context.Background()
	// Set this in app.yaml when running in production.
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	var err error
	datastoreClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer func(datastoreClient *datastore.Client) {
		err := datastoreClient.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(datastoreClient)

	http.HandleFunc("/invitations", BasicAuth(getAllInvit))
	http.HandleFunc("/create", BasicAuth(createInvit))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getAllInvit(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/invitations" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	// Get a list of the most recent visits.
	invitations, err := models.GetAllInvitations(ctx, datastoreClient)
	if err != nil {
		msg := fmt.Sprintf("Could not get recent visits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(invitations)
	_, err = w.Write(data)
	if err != nil {
		msg := fmt.Sprintf("Could not get invits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func createInvit(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	invitData := models.InvitData{}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&invitData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := models.CreateInvitation(ctx, datastoreClient, invitData.IsWedding, invitData.IsReception, invitData.HasCompanion)
	if err != nil {
		msg := fmt.Sprintf("Could not save visit: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	invit, err := models.GetInvitation(ctx, datastoreClient, key)
	if err != nil {
		msg := fmt.Sprintf("Could not save visit: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(invit)
	_, err = w.Write(data)
}

func unauthorised(rw http.ResponseWriter) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	rw.WriteHeader(http.StatusUnauthorized)

}

func BasicAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			log.Println("Error getting all credentials")
			unauthorised(w)
			return
		}
		ctx := context.Background()

		query := datastore.NewQuery("User").Filter("Username =", u)
		var users []*models.User

		if _, err := datastoreClient.GetAll(ctx, query, &users); err != nil {
			log.Println("Error getting all users", err)
			unauthorised(w)
			return

		}
		log.Println("THIS IS USERS:", users)

		if len(users) == 0 {
			log.Println("There are no users in users list")

			unauthorised(w)
			return
		}

		if p != users[0].Password {
			log.Println("Password didnt match")
			unauthorised(w)
			return
		}

		// log.Println("some usr:", u, "pass:", p, "and is ok?:", ok)
		f(w, r)
	}
}
