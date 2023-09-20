package server

import (
	"SimpleReader/server/api"
	"SimpleReader/server/api/utils"
	"SimpleReader/server/models/global"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
)

func Start() {
	gin.SetMode(gin.DebugMode)

	go utils.MakeCover()

	storage.NewStorage("db/storage")

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(storage.Sessions())

	r.GET("/back.jpg", func(c *gin.Context) {
		c.File("public/img/back.jpg")
	})

	r.Static("/js", "public/js")
	r.Static("/css", "public/css")
	r.Static("/img", "public/img")

	r.LoadHTMLGlob("public/views/*/*.gohtml")

	api.SetRoutes(r)

	r.Run(":8080")

	global.Stoped = true
}
