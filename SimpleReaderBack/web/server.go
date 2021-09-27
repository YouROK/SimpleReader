package web

import (
	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/api"
	"SimpleReader/web/settings"
)

func Start() {
	settings.Open()
	defer settings.Close()

	app := fiber.New(fiber.Config{BodyLimit: 10 * 1024 * 1024})

	app.Post("/api/login", api.Login)

	app.Get("/api/register/:id", api.Register)
	app.Post("/api/register/:id", api.Register)

	app.Post("/api/upload", api.Upload)

	app.Listen(":9000")
}
