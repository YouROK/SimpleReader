package api

import (
	"SimpleReader/server/api/pages"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	r.GET("/", pages.MainPage)
}
