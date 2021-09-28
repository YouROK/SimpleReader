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
		Login    string `json:"login"`
		PassHash string `json:"pass"`
	}{}

	if err := c.QueryParser(&payload); err != nil {
		time.Sleep(time.Second * 3)
		c.Status(http.StatusBadRequest)
		return c.SendString(err.Error())
	}

	if payload.Login == "" || payload.PassHash == "" {
		time.Sleep(time.Second * 3)
		c.Status(http.StatusBadRequest)
		return c.SendString("Заполните все поля")
	}

	ses := session.Get(c)
	userDb := storage.GetUser(payload.Login)
	if userDb == nil {
		time.Sleep(time.Second * 5)
		c.Status(http.StatusBadRequest)
		return c.SendString("Неправильное имя пользователя")
	}

	if userDb.PassHash != payload.PassHash {
		time.Sleep(time.Second * 5)
		c.Status(http.StatusBadRequest)
		return c.SendString("Неправильный пароль")
	}

	if userDb.Role < 0 {
		time.Sleep(time.Second * 5)
		c.Status(http.StatusBadRequest)
		return c.SendString("Вы забанены")
	}

	ses.Set("User", userDb)
	c.Status(200)
	return nil
}
