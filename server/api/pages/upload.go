package pages

import (
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadPage(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "upload.gohtml", nil)
}

func UploadBook(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.String(http.StatusUnauthorized, "Ошибка загрузки, нужно войти на сайт под своим именем")
		return
	}

	form, err := c.MultipartForm()
	if err == nil {
		if form == nil || len(form.File) == 0 {
			c.String(400, "Выберите файлы для загрузки")
			return
		} else {
			defer form.RemoveAll()
			if files, ok := form.File["fb2files"]; ok {
				loads := make([]string, 0)
				errloads := make([]string, 0)
				for _, f := range files {
					books, err := storage.GetBookStorage().WriteBook(f)
					if err != nil {
						errloads = append(errloads, f.Filename)
					} else {
						loads = append(loads, books...)
					}
				}
				for _, b := range loads {
					if _, ok := ses.User.ReadingBooks[b]; !ok {
						ses.User.ReadingBooks[b] = &user.BookInfo{BookHash: b}
					}
				}
				storage.GetUserStorage().SaveUser(ses.User)
				c.JSON(200, loads)
				return
			}
			c.String(400, "Ошибка загрузки книги")
			return
		}
	} else {
		c.String(500, "Ошибка загрузки, "+err.Error())
		return
	}
}
