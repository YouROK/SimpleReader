package pages

import (
	"SimpleReader/server/models/fb2"
	"SimpleReader/server/models/session"
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"sort"
)

func MainPage(c *gin.Context) {
	books := storage.GetSessionStorage().GetSession(c).User.ReadingBooks
	ses := storage.GetSession(c)

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

	if len(descs) < 9 {
		books := storage.GetBookStorage().GetBooks()
		rand.Shuffle(len(books), func(i, j int) { books[i], books[j] = books[j], books[i] })
		for _, b := range books {
			desc, _ := storage.GetBookStorage().GetBookDesc(b)
			if desc != nil {
				descs = append(descs, desc)
			}
			if len(descs) > 9 {
				break
			}
		}
	}

	c.HTML(http.StatusOK, "main.gohtml", struct {
		Books   []*fb2.XMLTitleInfo
		Session *session.Session
	}{Books: descs, Session: ses})

}
