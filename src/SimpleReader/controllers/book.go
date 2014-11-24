package controllers

import (
	"SimpleReader/models/fbreader"
	"SimpleReader/models/sessions"
	"SimpleReader/models/storage"
	"SimpleReader/models/users"
	"SimpleReader/server"
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type Book struct {
	server.Controller
}

func (book *Book) InitController(serv *server.Server) {
	serv.GetMartini().Router.Get("/book/:book", book.GetBookInfo)
	serv.GetMartini().Router.Get("/book/:book/read", book.GetBookReadPage)
	serv.GetMartini().Router.Get("/book/:book/read/:page", book.SetBookReadPage)
	serv.GetMartini().Router.Get("/getcover/:book", book.GetCover)
	serv.GetMartini().Router.Get("/getimage/:book/:id", book.GetImage)
	serv.GetMartini().Router.Post("/getcontent", book.GetContent)
	serv.GetMartini().Router.Post("/book/:book/:page", book.SetPage)
	serv.GetMartini().Router.Get("/lastbook", book.GetLastBook)
	serv.GetMartini().Router.Post("/getchapters/:book", book.GetChapters)
}

func (b *Book) GetBookReadPage(stor storage.Storages, req *http.Request, r render.Render, params martini.Params) {
	bookHash := params["book"]
	if bookHash != "" {
		type resp struct {
			BookFb   *fbreader.FBParser
			BookInfo *users.BookInfo
			Session  *sessions.Session
		}
		book := stor.GetBookStorage().GetBook(bookHash)
		if book == nil {
			r.Error(404)
			return
		}
		ses := stor.GetSessionStorage().GetSession(req)
		bookInfo := ses.User.GetBookInfo(bookHash)
		if bookInfo.SetPage {
			defer func() {
				bookInfo.SetPage = false
				ses.User.SetBookInfo(bookInfo)
			}()
		}
		r.HTML(200, "bookpages/bookpage", &resp{book, &bookInfo, ses})
	} else {
		r.Error(404)
	}
}

func (b *Book) SetBookReadPage(stor storage.Storages, req *http.Request, res http.ResponseWriter, r render.Render, params martini.Params) {
	bookHash := params["book"]
	pageStr := params["page"]
	page, err := strconv.Atoi(pageStr)
	if bookHash != "" && err == nil {
		ses := stor.GetSessionStorage().GetSession(req)
		bookInfo := ses.User.GetBookInfo(bookHash)

		bookInfo.LastRead = time.Now().UTC()
		bookInfo.LastReadPage = page
		bookInfo.SetPage = true
		ses.User.SetBookInfo(bookInfo)
		stor.GetUserStorage().SaveUser(ses.User)
		r.Redirect("/book/" + bookHash + "/read")
		return
	}
	r.Error(404)
}

func (b *Book) GetBookInfo(stor storage.Storages, req *http.Request, r render.Render, params martini.Params) {
	bookHash := params["book"]
	if bookHash != "" {
		book := stor.GetBookStorage().GetBook(bookHash)
		type BookInfoPage struct {
			Book    *fbreader.FBParser
			Session *sessions.Session
		}
		r.HTML(200, "bookpages/bookinfo", BookInfoPage{book, stor.GetSessionStorage().GetSession(req)})
	}
}

func (b *Book) GetCover(stor storage.Storages, res http.ResponseWriter, req *http.Request, params martini.Params) {
	bookHash := params["book"]
	if bookHash != "" {
		bookdesc, err := stor.GetBookStorage().GetBookDesc(bookHash)
		if err == nil && bookdesc.Coverpage.ImgLink != "" {
			coverPath := path.Join(stor.GetBookStorage().GetPath(), bookHash, bookdesc.Coverpage.ImgLink[1:])
			if s, e := os.Stat(coverPath); e == nil && s.Mode().IsRegular() {
				http.ServeFile(res, req, coverPath)
				return
			}
		}
	}
	http.ServeFile(res, req, "public/images/bookcover.jpg")
}

func (b *Book) GetImage(stor storage.Storages, res http.ResponseWriter, req *http.Request, params martini.Params) {
	bookHash := params["book"]
	imgId := params["id"]
	if bookHash != "" && imgId != "" {
		imgs, err := stor.GetBookStorage().GetBookImgs(bookHash)
		if err == nil {
			for _, img := range imgs {
				if imgId == img {
					imgPath := path.Join(stor.GetBookStorage().GetPath(), bookHash, imgId)
					http.ServeFile(res, req, imgPath)
					return
				}
			}
		}
	}
	http.NotFound(res, req)
}

func (b *Book) GetContent(stor storage.Storages, res http.ResponseWriter, req *http.Request) {
	type jsonreq struct {
		BookId string
		Start  int
		Count  int
	}
	jreq := jsonreq{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&jreq)
	if err == nil {
		book := stor.GetBookStorage().GetBook(jreq.BookId)
		if book == nil {
			http.NotFound(res, req)
			return
		}
		if jreq.Start <= book.GetContentCount() {
			end := jreq.Start + jreq.Count
			if end > book.GetContentCount() {
				end = book.GetContentCount()
			}
			content := (*book.GetContent())[jreq.Start:end]
			enc := json.NewEncoder(res)
			enc.Encode(content)
			res.Header().Set("Content-Type", "application/json;charset=utf-8")
			return
		} else {
			http.Error(res, "Book content length is less then request", 400)
		}
	} else {
		log.Println("Error decode json (getcontent)")
	}

	http.Error(res, err.Error(), 500)
}

func (b *Book) SetPage(stor storage.Storages, req *http.Request, params martini.Params) (int, string) {
	bookHash := params["book"]
	pageStr := params["page"]

	bookPrs := stor.GetBookStorage().GetBook(bookHash)

	if bookHash == "" || bookPrs == nil {
		return 404, "Book with id '" + bookHash + "' not found"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 400, err.Error()
	}
	if page < 0 || page > bookPrs.GetContentCount() {
		return 400, "Wrong page number " + pageStr
	}

	ses := stor.GetSessionStorage().GetSession(req)
	bookInfo := ses.User.GetBookInfo(bookHash)

	bookInfo.LastRead = time.Now().UTC()
	bookInfo.LastReadPage = page
	ses.User.SetBookInfo(bookInfo)
	stor.GetUserStorage().SaveUser(ses.User)

	if ses.Login != 1 {
		return 401, ""
	}

	return 200, ""
}

func (b *Book) GetLastBook(stor storage.Storages, req *http.Request, r render.Render) {
	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login != 1 {
		r.Redirect("/login")
		return
	}

	var book users.BookInfo
	for _, v := range ses.User.ReadingBooks {
		if book.BookHash == "" {
			book = v
			continue
		}
		if v.LastRead.After(book.LastRead) {
			book = v
		}
	}

	if book.BookHash == "" {
		r.Redirect("/library")
		return
	}

	r.Redirect("/book/" + book.BookHash + "/read")
}

func (b *Book) GetChapters(stor storage.Storages, req *http.Request, res http.ResponseWriter, params martini.Params) (int, string) {
	bookHash := params["book"]

	bookPrs := stor.GetBookStorage().GetBook(bookHash)
	res.Header().Set("Content-Type", "application/json;charset=utf-8")

	if bookHash == "" || bookPrs == nil {
		return 404, getMsg("Book with id '" + bookHash + "' not found")
	}

	chapters := bookPrs.GetChapters()
	if chapters == nil {
		return 500, getMsg("Error get content of book with id '" + bookHash + "'")
	}

	bJson, err := json.Marshal(*chapters)
	if err != nil {
		return 500, getMsg("Error parse chapters of book with id '" + bookHash + "'")
	}

	return 200, string(bJson)
}
