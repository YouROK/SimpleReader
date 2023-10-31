package pages

import (
	"SimpleReader/server/api/fbreader"
	"SimpleReader/server/models/session"
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func BookReadPage(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	bookHash := c.Param("hash")
	if bookHash != "" {
		type resp struct {
			BookFb   *fbreader.FBParser
			BookInfo *user.BookInfo
			Session  *session.Session
		}
		book := storage.GetBookStorage().GetBook(bookHash)
		if book == nil {
			c.Status(404)
			return
		}
		ses := storage.GetSessionStorage().GetSession(c)
		bookInfo := ses.User.GetBookInfo(bookHash)
		c.HTML(200, "bookread.gohtml", &resp{book, bookInfo, ses})
	} else {
		c.Status(404)
	}
}

func BookReadPageSet(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	bookHash := c.Param("hash")
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if bookHash != "" && err == nil {
		bookInfo := ses.User.GetBookInfo(bookHash)

		bookInfo.LastRead = time.Now().UTC()
		bookInfo.LastReadPage = page
		bookInfo.SetPage = true
		ses.User.SetBookInfo(bookInfo)
		storage.GetUserStorage().SaveUser(ses.User)
		c.Redirect(http.StatusFound, "/book/"+bookHash+"/read")
		return
	}
	c.Status(404)
}

func BookReadSet(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.String(http.StatusUnauthorized, "Вы не авторизованы")
		return
	}

	bookHash := c.Param("hash")
	pageStr := c.Param("page")

	bookPrs := storage.GetBookStorage().GetBook(bookHash)

	if bookHash == "" || bookPrs == nil {
		c.String(404, "Book with id '"+bookHash+"' not found")
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	if page < 0 || page > bookPrs.GetContentCount() {
		c.String(400, "Wrong page number "+pageStr)
		return
	}

	bookInfo := ses.User.GetBookInfo(bookHash)

	bookInfo.LastRead = time.Now().UTC()
	bookInfo.LastReadPage = page
	ses.User.SetBookInfo(bookInfo)
	err = storage.GetUserStorage().SaveUser(ses.User)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(200)
}

func GetContent(c *gin.Context) {
	ses := storage.GetSession(c)
	if !ses.IsLogin() {
		c.String(http.StatusUnauthorized, "Вы не авторизованы")
		return
	}

	type jsonreq struct {
		BookId string
		Start  int
		Count  int
	}
	jreq := jsonreq{}
	err := c.BindJSON(&jreq)
	if err == nil {
		book := storage.GetBookStorage().GetBook(jreq.BookId)
		if book == nil {
			c.Status(404)
			return
		}
		if jreq.Start <= book.GetContentCount() {
			end := jreq.Start + jreq.Count
			if end > book.GetContentCount() {
				end = book.GetContentCount()
			}
			content := book.GetContent()[jreq.Start:end]
			c.JSON(200, content)
			return
		} else {
			c.String(400, "Book content length is less then request")
		}
	} else {
		log.Println("Error decode json (getcontent):", err)
	}

	c.String(500, err.Error())
}
