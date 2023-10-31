package pages

import (
	"SimpleReader/server/models/fb2"
	"SimpleReader/server/models/session"
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sort"
	"strings"
)

func MyBooksPage(c *gin.Context) {
	ses := storage.GetSession(c)

	if !ses.IsLogin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	books := storage.GetSessionStorage().GetSession(c).User.ReadingBooks

	var descs []*fb2.XMLTitleInfo
	if len(books) > 0 {
		var list []*user.BookInfo
		for _, info := range books {
			list = append(list, info)
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].LastRead.After(list[j].LastRead)
		})
		for _, b := range list {
			desc, _ := storage.GetBookStorage().GetBookDesc(b.BookHash)
			if desc != nil {
				descs = append(descs, desc)
			}
		}
	}

	c.HTML(http.StatusOK, "mybooks.gohtml", struct {
		Books   []*fb2.XMLTitleInfo
		Session *session.Session
	}{Books: descs, Session: ses})

}

func RemBook(c *gin.Context) {
	ses := storage.GetSession(c)

	if !ses.IsLogin() {
		c.Status(http.StatusUnauthorized)
		return
	}

	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	hash := strings.TrimSpace(strings.ToLower(string(buf)))
	if _, ok := ses.User.ReadingBooks[hash]; ok {
		delete(ses.User.ReadingBooks, hash)
		storage.GetUserStorage().SaveUser(ses.User)
		c.Status(200)
		return
	}
	c.String(http.StatusBadRequest, "Книга не найдена")
	return
}
