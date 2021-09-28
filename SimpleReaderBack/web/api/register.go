package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/models"
	"SimpleReader/web/settings"
)

func Register(c *fiber.Ctx) error {
	hash := c.Params("hash", "")
	if hash == "" {
		c.Status(http.StatusNotFound)
		return nil
	}
	regLink := settings.GetLink(hash)
	if regLink == nil {
		c.Status(http.StatusNotFound)
		return nil
	}
	return c.JSON(regLink)
}

func RegisterSetData(c *fiber.Ctx) error {
	hash := c.Params("hash", "")
	if hash == "" {
		c.Status(http.StatusBadRequest)
		return nil
	}
	regLink := settings.GetLink(hash)
	if regLink == nil {
		c.Status(http.StatusBadRequest)
		return nil
	}

	usr := new(models.User)
	err := c.QueryParser(&usr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return err
	}

	return c.JSON(regLink)
}
