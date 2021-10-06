package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/models"
	"SimpleReader/web/session"
)

func GetStyle(c *fiber.Ctx) error {
	ses := session.Get(c)
	usr, _ := ses.Get("User").(*models.User)
	if usr == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	lst := make([]models.BookInfo, 0)
	for _, info := range usr.ReadBooks {
		lst = append(lst, info)
	}

	return c.JSON(lst)
}

func GetReadBooks(c *fiber.Ctx) error {
	ses := session.Get(c)
	usr, _ := ses.Get("User").(*models.User)
	if usr == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	lst := make([]models.BookInfo, 0)
	for _, info := range usr.ReadBooks {
		lst = append(lst, info)
	}

	return c.JSON(lst)
}
