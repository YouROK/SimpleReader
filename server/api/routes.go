package api

import (
	"SimpleReader/server/api/pages"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	r.GET("/", pages.MainPage)

	r.GET("/login", pages.LoginPage)
	r.POST("/login", pages.Login)
	r.GET("/logout", pages.Logout)
	r.GET("/getkey", pages.GetKey)

	r.GET("/mybooks", pages.MyBooksPage)
	r.POST("/rembook", pages.RemBook)

	r.GET("/upload", pages.UploadPage)
	r.POST("/upload", pages.UploadBook)

	r.GET("/settings", pages.SetsPage)
	r.POST("/settings", pages.SetStyle)

	r.GET("/search", pages.SearchPage)
	r.POST("/search", pages.SearchBooks)

	r.GET("/authors", pages.SearchAuthorPage)
	r.POST("/authors", pages.SearchAuthors)
	r.GET("/author/:hash", pages.AuthorPage)

	r.GET("/book/:hash", pages.DescPage)
	r.GET("/cover/:hash", pages.Cover)
	r.GET("/lastbook", pages.LastBook)

	r.GET("/book/:hash/read", pages.BookReadPage)
	r.GET("/book/:hash/read/:page", pages.BookReadPage)
	r.POST("/book/:hash/:page", pages.BookReadSet)
	r.POST("/getcontent", pages.GetContent)

	r.GET("/back.png", func(c *gin.Context) {
		c.File("public/img/back.png")
	})
}
