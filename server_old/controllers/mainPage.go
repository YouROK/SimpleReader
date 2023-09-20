package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/martini-contrib/render"

	"SimpleReader/server_old/models/fbreader"
	"SimpleReader/server_old/models/sessions"
	"SimpleReader/server_old/models/storage"
	"SimpleReader/server_old/models/users"
	"SimpleReader/server_old/server"
)

type MainPage struct {
	server.Controller
}

func (h *MainPage) InitController(serv *server.Server) {
	serv.GetMartini().Router.Get("/", h.Handler)
}

func (h *MainPage) Handler(stor storage.Storages, req *http.Request, r render.Render) {
	booksHash := stor.GetBookStorage().GetBooks()
	BookCount := 5
	rand.Seed(time.Now().Unix())

	bookslist := make([]*fbreader.XMLTitleInfo, 0, BookCount+1)
	readingBooks := stor.GetSessionStorage().GetSession(req).User.ReadingBooks
	if len(readingBooks) > 0 {
		var book users.BookInfo
		for _, v := range readingBooks {
			if book.BookHash == "" {
				book = v
				continue
			}
			if v.LastRead.After(book.LastRead) {
				book = v
			}
		}
		if book.BookHash != "" {
			desc, err := storage.GetStorage().GetBookStorage().GetBookDesc(book.BookHash)
			if err == nil {
				desc.Hash = book.BookHash
				bookslist = append(bookslist, desc)
				BookCount += 1
			}
		}
	}

	i := 0
	if len(booksHash) > 0 {
		for len(bookslist) < BookCount {
			hash := booksHash[rand.Intn(len(booksHash))]
			desc, err := storage.GetStorage().GetBookStorage().GetBookDesc(hash)
			if err == nil {
				desc.Hash = hash
				bookslist = append(bookslist, desc)
			}
			i++
			if i > 9 {
				break
			}
		}
	}

	type MainPage struct {
		Books   []*fbreader.XMLTitleInfo
		Session *sessions.Session
	}
	r.HTML(200, "mainpage/mainpage", MainPage{bookslist, stor.GetSessionStorage().GetSession(req)})
}
