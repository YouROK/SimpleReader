package api

import (
	"io/ioutil"
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

	path := filepath.Join(settings.Path, "books", hash, "book.fb2")
	return c.SendFile(path, true)
}

func GetBooks(c *fiber.Ctx) error {
	ses := session.Get(c)
	usr, _ := ses.Get("User").(*models.User)
	if usr == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	path := filepath.Join(settings.Path, "books")
	ff, err := ioutil.ReadDir(path)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	lst := make([]string, 0)
	for _, info := range ff {
		if info.IsDir() {
			lst = append(lst, info.Name())
		}
	}
	return c.JSON(lst)
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

	path := filepath.Join(settings.Path, "books", hash, name)
	return c.SendFile(path)
}

func GetBookDesc(c *fiber.Ctx) error {
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
	buf, err := ioutil.ReadFile(filepath.Join(settings.Path, "books", hash, "desc.json"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}
	c.Response().Header.SetContentType("application/json;charset=UTF-8")
	return c.SendString(string(buf))
}
