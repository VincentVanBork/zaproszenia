// Sample datastore demonstrates use of the cloud.google.com/go/datastore package from App Engine flexible.
package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"log"
	"os"
	"zaproszenia/controllers"
	"zaproszenia/models"
	"zaproszenia/utils"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	engine := html.New("./static", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use("/admin", basicauth.New(basicauth.Config{
		Users: utils.GetUsers(datastoreClient),
	}))
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	app.Get("/:key", func(c *fiber.Ctx) error {
		Token := c.Query("token", "")
		KeyString := c.Params("key")
		Key, KeyEncodeErr := datastore.DecodeKey(KeyString)
		if KeyEncodeErr != nil {
			return err
		}
		invitation, Inviterr := models.GetInvitation(ctx, datastoreClient, Key)
		if Inviterr != nil {
			return Inviterr
		}
		if invitation.Token == Token {
			return c.Render("index", fiber.Map{})
		}
		return errors.New("Could not match token to your invitation")
	})

	ControllerInvit := controllers.InvitationsController{Objects: datastoreClient}
	app.Get("/api/invitations/:key", ControllerInvit.GetFullInvitation)
	app.Get("/admin/invitations/:key", ControllerInvit.GetOne)
	app.Get("/admin/invitations", ControllerInvit.GetAll)
	app.Post("/admin/invitations", ControllerInvit.Create)

	ControllerGuest := controllers.GuestController{Objects: datastoreClient}
	app.Get("/admin/guests", ControllerGuest.GetAll)
	app.Post("/admin/guests", ControllerGuest.Create)

	log.Fatal(app.Listen(":" + port))
}
