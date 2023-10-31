package pages

import (
	"SimpleReader/server/api/fbreader"
	"SimpleReader/server/models/session"
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func DescPage(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.Status(404)
		return
	}

	book := storage.GetBookStorage().GetBook(hash)
	if book == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	ses := storage.GetSession(c)

	c.HTML(http.StatusOK, "desc.gohtml", struct {
		Book    *fbreader.FBParser
		Session *session.Session
	}{Book: book, Session: ses})
}

func Cover(c *gin.Context) {
	hash := c.Param("hash")
	if hash != "" {
		desc, err := storage.GetBookStorage().GetBookDesc(hash)
		if err == nil && desc.Coverpage.ImgLink != "" {
			coverPath := path.Join(storage.GetBookStorage().GetPath(), hash, desc.Coverpage.ImgLink[1:])
			if s, e := os.Stat(coverPath); e == nil && s.Mode().IsRegular() {
				c.File(coverPath)
				return
			}
		}
	}
	c.File("public/img/bookcover.jpg")
}

func LastBook(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var book *user.BookInfo
	for _, v := range ses.User.ReadingBooks {
		if book == nil || book.BookHash == "" {
			book = v
			continue
		}
		if v.LastRead.After(book.LastRead) {
			book = v
		}
	}

	if book != nil && book.BookHash == "" {
		c.Redirect(http.StatusFound, "/mybooks")
		return
	}

	c.Redirect(http.StatusFound, "/book/"+book.BookHash+"/read")
}
