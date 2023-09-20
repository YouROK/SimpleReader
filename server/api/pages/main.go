package pages

import (
	"SimpleReader/server/models/fb2"
	"SimpleReader/server/models/session"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MainPage(c *gin.Context) {
	//books := storage.GetSessionStorage().GetSession(c).User.ReadingBooks
	books := storage.GetBookStorage().GetBooks()
	ses := storage.GetSession(c)

	var descs []*fb2.XMLTitleInfo

	for _, b := range books {
		desc, _ := storage.GetBookStorage().GetBookDesc(b)
		if desc != nil {
			descs = append(descs, desc)
		}
	}

	c.HTML(http.StatusOK, "main/main.gohtml", struct {
		Books   []*fb2.XMLTitleInfo
		Session *session.Session
	}{Books: descs, Session: ses})

}
