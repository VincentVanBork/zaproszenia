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
	"zaproszenia/utils"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
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

	http.HandleFunc("/", basicAuth(getAllInvit))
	http.HandleFunc("/create", basicAuth(createInvit))
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

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for the provided and expected
			// usernames and passwords.
			passwordMatch := utils.CheckPasswordHash(password, os.Getenv("AUTH_PASSWORD"))
			usernameMatch := utils.CheckPasswordHash(username, os.Getenv("AUTH_USERNAME"))

			// If the username and password are correct, then call
			// the next handler in the chain. Make sure to return
			// afterwards, so that none of the code below is run.
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
