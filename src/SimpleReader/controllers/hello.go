package controllers

import (
	"SimpleReader/models/fbreader"
	"SimpleReader/models/storage"
	"SimpleReader/server"
	"github.com/martini-contrib/render"
	"log"
	"math/rand"
	"net/http"
)

type Hello struct {
	server.Controller
	server *server.Server
	//fb     *fbreader.FBParser
}

func (h *Hello) InitController(serv *server.Server) {
	//serv.GetMartini().Router.Get("/hello", h.Handler)
	//h.fb = fbreader.GetParser("db/storage/books/2c6e9f99c12b3d4104f826ae3c52455f/book.fb2")
	//h.fb.ParseFB()
	//h.server = serv
}

func (h *Hello) Handler(stor storage.Storages, req *http.Request, r render.Render) {
	booksHash := stor.GetBookStorage().GetBooks()
	log.Println(booksHash)
	books := make([]*fbreader.FBParser, 5)
	for i := 0; i < 5; i++ {
		hash := booksHash[rand.Intn(len(booksHash))]
		books[i] = stor.GetBookStorage().GetBook(hash)
	}
	r.HTML(200, "mainpage/mainpage", books, render.HTMLOptions{""})
}
