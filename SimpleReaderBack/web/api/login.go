package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/session"
	"SimpleReader/web/storage"
)

func IsLogin(c *fiber.Ctx) error {
	ses := session.Get(c)
	usr := ses.Get("User")
	if usr == nil {
		c.Status(http.StatusUnauthorized)
	}
	return nil
}

func Login(c *fiber.Ctx) error {
	payload := struct {
		EMail    string `json:"email"`
		PassHash string `json:"pass"`
	}{}

	err := json.Unmarshal(c.Body(), &payload)
	if err != nil {
		time.Sleep(time.Second * 3)
		c.Status(http.StatusBadRequest)
		return c.SendString(err.Error())
	}

	if payload.EMail == "" || payload.PassHash == "" {
		time.Sleep(time.Second * 3)
		c.Status(http.StatusBadRequest)
		return c.SendString("Заполните все поля")
	}

	user := storage.GetUser(payload.EMail)
	if user == nil {
		time.Sleep(time.Second * 5)
		c.Status(http.StatusBadRequest)
		return c.SendString("Неправильное имя пользователя")
	}

	if user.PassHash != payload.PassHash {
		time.Sleep(time.Second * 5)
		c.Status(http.StatusBadRequest)
		return c.SendString("Неправильный пароль")
	}

	if user.Role < 0 {
		time.Sleep(time.Second * 5)
		c.Status(http.StatusBadRequest)
		return c.SendString("Вы забанены")
	}
	ses := session.Get(c)
	ses.Set("User", user)
	c.Status(200)
	// TODO понять что не так с cookies
	return c.SendString(ses.ID())
}
