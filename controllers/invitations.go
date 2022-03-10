package controllers

import (
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
