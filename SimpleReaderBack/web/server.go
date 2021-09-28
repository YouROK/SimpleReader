package web

import (
	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/api"
	"SimpleReader/web/settings"
)

func Start() {
	// settings.Path = filepath.Dir(os.Args[0])
	settings.Path = "db"

	settings.Open()
	defer settings.Close()

	app := fiber.New(fiber.Config{BodyLimit: 10 * 1024 * 1024})

	app.Post("/api/login", api.Login)

	app.Get("/api/register/:hash", api.Register)
	app.Post("/api/register/:hash", api.RegisterSetData)

	app.Post("/api/upload", api.Upload)

	app.Get("/api/reads", api.GetReadBooks)

	app.Get("/api/book/:hash", api.GetBook)
	app.Get("/api/bin/:hash/:name", api.GetBin)

	app.Listen(":9000")
}
