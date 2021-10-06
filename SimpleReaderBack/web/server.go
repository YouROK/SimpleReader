package web

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"SimpleReader/web/api"
	"SimpleReader/web/session"
	"SimpleReader/web/settings"
	"SimpleReader/web/utils"
)

func Start() {
	// settings.Path = filepath.Dir(os.Args[0])
	settings.Path = "db"
	os.MkdirAll(settings.Path, 0777)

	settings.Open()
	settings.InitLinksChecker()
	defer settings.Close()

	go utils.MakeCover()

	app := fiber.New(fiber.Config{BodyLimit: 50 * 1024 * 1024})
	session.Init()
	// TODO remove on release
	app.Use(cors.New())

	app.Get("/api/login", api.IsLogin)
	app.Post("/api/login", api.Login)

	app.Get("/api/register/:hash", api.RegisterGetEmail)
	app.Post("/api/register/:hash", api.RegisterSetData)

	app.Get("/api/user/reads", api.GetReadBooks)
	app.Get("/api/user/style", api.GetStyle)

	app.Post("/api/book/upload", api.Upload)
	app.Get("/api/book/get/:hash", api.GetBook)
	app.Get("/api/book/all", api.GetBooks)
	app.Get("/api/book/bin/:hash/:name", api.GetBin)
	app.Get("/api/book/desc/:hash", api.GetBookDesc)

	app.Get("/api/search", api.Search)

	app.Get("/img/back.jpg", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(settings.Path, "img", "back.jpg"))
	})
	app.Get("/img/cover.jpg", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(settings.Path, "img", "cover.jpg"))
	})

	app.Listen(":9000")
}
