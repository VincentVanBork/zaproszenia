// Sample datastore demonstrates use of the cloud.google.com/go/datastore package from App Engine flexible.
package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"os"
	"zaproszenia/models"
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

	http.HandleFunc("/", getAllInvit)
	http.HandleFunc("/create", createInvit)
	appengine.Main()
}

func getAllInvit(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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
