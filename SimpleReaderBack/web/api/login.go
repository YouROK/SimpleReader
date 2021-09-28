package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/session"
	"SimpleReader/web/storage"
)

func Login(c *fiber.Ctx) error {
	payload := struct {
		EMail    string `json:"email"`
		PassHash string `json:"pass"`
	}{}

	if err := c.QueryParser(&payload); err != nil {
		time.Sleep(time.Second * 3)
		c.Status(http.StatusBadRequest)
		return c.SendString(err.Error())
	}

	if payload.EMail == "" || payload.PassHash == "" {
		time.Sleep(time.Second * 3)
		c.Status(http.StatusBadRequest)
		return c.SendString("Заполните все поля")
	}

	ses := session.Get(c)
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

	ses.Set("User", user)
	c.Status(200)
	return nil
}
