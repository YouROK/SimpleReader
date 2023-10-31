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
	go ParseDir()

	r := gin.Default()

	r.MaxMultipartMemory = 32 << 20

	r.Use(gin.Recovery())
	r.Use(storage.Sessions())

	r.Static("/js", "public/js")
	r.Static("/css", "public/css")
	r.Static("/img", "public/img")

	r.LoadHTMLGlob("public/views/*.gohtml")

	api.SetRoutes(r)

	r.Run(":8080")

	global.Stoped = true
}
