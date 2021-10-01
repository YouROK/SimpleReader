package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/models"
	"SimpleReader/web/settings"
	"SimpleReader/web/storage"
)

func RegisterGetEmail(c *fiber.Ctx) error {
	hash := c.Params("hash", "")
	if hash == "" {
		c.Status(http.StatusNotFound)
		return nil
	}
	regLink := settings.GetLink(hash)
	if regLink == nil {
		time.Sleep(5 * time.Second) // wait for bruteforce
		c.Status(http.StatusNotFound)
		return nil
	}
	return c.SendString(regLink.EMail)
}

func RegisterSetData(c *fiber.Ctx) error {
	hash := c.Params("hash", "")
	if hash == "" {
		c.Status(http.StatusBadRequest)
		return nil
	}
	regLink := settings.GetLink(hash)
	if regLink == nil {
		time.Sleep(5 * time.Second) // wait for bruteforce
		c.Status(http.StatusBadRequest)
		return errors.New("Произошла ошибка, обратитесь к создателю")
	}

	payload := struct {
		Login    string `json:"login"`
		PassHash string `json:"pass"`
	}{}

	err := json.Unmarshal(c.Body(), &payload)
	if err != nil || payload.Login == "" || payload.PassHash == "" {
		c.Status(http.StatusBadRequest)
		return err
	}

	usrDB := storage.GetUser(regLink.EMail)
	if usrDB == nil {
		usrDB = new(models.User)
	}
	usrDB.Email = regLink.EMail
	usrDB.Login = payload.Login
	usrDB.PassHash = payload.PassHash
	storage.SetUser(usrDB)
	settings.RemLink(hash)
	c.Status(http.StatusOK)
	return nil
}
