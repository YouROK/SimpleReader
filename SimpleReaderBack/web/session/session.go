package session

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	store *session.Store
)

func init() {
	cfg := session.ConfigDefault
	cfg.Expiration = 48 * time.Hour
	store = session.New(cfg)
}

func Get(c *fiber.Ctx) *session.Session {
	stor, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	return stor
}
