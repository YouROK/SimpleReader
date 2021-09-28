package api

import (
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/models"
	"SimpleReader/web/session"
	"SimpleReader/web/settings"
)

func GetBook(c *fiber.Ctx) error {
	ses := session.Get(c)
	usr, _ := ses.Get("User").(*models.User)
	if usr == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	hash := c.Params("hash", "")
	if hash == "" {
		c.Status(http.StatusBadRequest)
		return nil
	}

	path := filepath.Join(settings.Path, hash, "book.fb2")
	return c.SendFile(path, true)
}

func GetBin(c *fiber.Ctx) error {
	ses := session.Get(c)
	usr, _ := ses.Get("User").(*models.User)
	if usr == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	hash := c.Params("hash", "")
	name := c.Params("name", "")
	if hash == "" || name == "" || name == "book.fb2" || name == "info.json" {
		c.Status(http.StatusBadRequest)
		return nil
	}

	path := filepath.Join(settings.Path, hash, name)
	return c.SendFile(path)
}
