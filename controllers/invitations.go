package controllers

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"zaproszenia/models"
)

type InvitationsController Controller

func (u *InvitationsController) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()
	invitations, err := models.GetAllInvitations(ctx, u.Objects)
	if err != nil {
		log.Println("Error getting invitations", err)
		return err
	}
	return c.JSON(invitations)
}

func (u *InvitationsController) Create(c *fiber.Ctx) error {
	ctx := context.Background()

	invitData := models.InvitData{}
	if err := c.BodyParser(&invitData); err != nil {
		log.Println("Error create invit BAD RQEUEST", err)
		return err
	}
	key, err := models.CreateInvitation(ctx, u.Objects, invitData.IsWedding, invitData.IsReception, invitData.HasCompanion)
	if err != nil {
		log.Println("Could not save visit: ", err)
		return err
	}
	invit, err := models.GetInvitation(ctx, u.Objects, key)
	return c.JSON(invit)
}

func (u *InvitationsController) GetFullInvitation(c *fiber.Ctx) error {
	ctx := context.Background()

	InvitationKey, keyerr := datastore.DecodeKey(c.Params("key"))
	if keyerr != nil {
		return keyerr
	}

	invit, err := models.GetInvitation(ctx, u.Objects, InvitationKey)
	if err != nil {
		return err
	}

	guests, err := models.GetGuestsByInvitationKey(ctx, u.Objects, InvitationKey)
	if err != nil {
		return err
	}

	return c.JSON(models.FullInvitationData{Invitation: invit, Guests: guests})
}
func (u *InvitationsController) GetOne(c *fiber.Ctx) error {
	ctx := context.Background()

	InvitationKey, keyerr := datastore.DecodeKey(c.Params("key"))
	if keyerr != nil {
		return keyerr
	}
	invit, err := models.GetInvitation(ctx, u.Objects, InvitationKey)
	if err != nil {
		return err
	}
	return c.JSON(invit)
}
