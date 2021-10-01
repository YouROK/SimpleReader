package web

import (
	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/api"
	"SimpleReader/web/session"
	"SimpleReader/web/settings"
	"SimpleReader/web/utils"
)

func Start() {
	// settings.Path = filepath.Dir(os.Args[0])
	settings.Path = "db"

	settings.Open()
	settings.InitLinksChecker()
	defer settings.Close()

	go utils.MakeCover()

	app := fiber.New(fiber.Config{BodyLimit: 10 * 1024 * 1024})
	session.Init()

	app.Get("/api/login", api.IsLogin)
	app.Post("/api/login", api.Login)

	app.Get("/api/register/:hash", api.RegisterGetEmail)
	app.Post("/api/register/:hash", api.RegisterSetData)

	app.Post("/api/upload", api.Upload)

	app.Get("/api/user/reads", api.GetReadBooks)
	app.Get("/api/user/style", api.GetStyle)

	app.Get("/api/book/:hash", api.GetBook)
	app.Get("/api/books", api.GetBooks)
	app.Get("/api/bin/:hash/:name", api.GetBin)

	app.Listen(":9000")
}
