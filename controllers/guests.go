package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"zaproszenia/models"
)

type GuestController Controller

func (u *GuestController) GetAll(c *fiber.Ctx) error {
	ctx := context.Background()
	guests, err := models.GetAllGuests(ctx, u.Objects)
	if err != nil {
		log.Println("Error getting guests", err)
		return err
	}
	return c.JSON(guests)
}

func (u *GuestController) Create(c *fiber.Ctx) error {
	ctx := context.Background()

	guestData := models.GuestData{}
	if parseerr := c.BodyParser(&guestData); parseerr != nil {
		log.Println("Error create invit BAD RQEUEST", parseerr)
		return parseerr
	}
	key, err := models.AssignGuest(ctx, u.Objects, guestData)
	if err != nil {
		log.Println("Could not save guest: ", err)
		return err
	}
	guest, guesterr := models.GetGuest(ctx, u.Objects, key)

	if guesterr != nil {
		return guesterr
	}

	return c.JSON(guest)
}
