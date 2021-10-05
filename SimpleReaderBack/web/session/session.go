package session

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"SimpleReader/web/models"
)

var (
	store *session.Store
)

func Init() {
	cfg := session.ConfigDefault
	cfg.Expiration = 48 * time.Hour
	// TODO remove
	cfg.CookieSameSite = "none"
	store = session.New(cfg)

	store.RegisterType(&models.User{})
}

func Get(c *fiber.Ctx) *session.Session {
	stor, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	return stor
}
