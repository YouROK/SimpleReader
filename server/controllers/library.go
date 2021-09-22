package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"SimpleReader/server/models/fbreader"
	"SimpleReader/server/models/sessions"
	"SimpleReader/server/models/storage"
	"SimpleReader/server/server"
)

type Library struct {
	server.Controller
}

func (l *Library) InitController(serv *server.Server) {
	serv.GetMartini().Router.Get("/search", l.SearchPage)
	serv.GetMartini().Router.Post("/search", l.SearchBooks)
	serv.GetMartini().Router.Get("/:user/books", l.UserBooks)
	serv.GetMartini().Router.Post("/removebooks", l.UserRemoveBooks)
	serv.GetMartini().Router.Get("/authors", l.AuthorsPage)
	serv.GetMartini().Router.Post("/authors", l.AuthorsSearch)
	serv.GetMartini().Router.Get("/author/:author", l.AuthorPage)
}

func (l *Library) SearchPage(stor storage.Storages, req *http.Request, r render.Render) {
	r.HTML(200, "library/searchPage", stor.GetSessionStorage().GetSession(req))
}

func (l *Library) SearchBooks(stor storage.Storages, req *http.Request, res http.ResponseWriter) (int, string) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return 500, err.Error()
	}

	srchStr := strings.ToLower(string(buf))
	if len(srchStr) < 3 {
		return 400, "Request is small"
	}

	booksHash := stor.GetBookStorage().GetBooks()
	const BookCount = 5
	doneParse := make(chan string, BookCount)
	booksmap := make(map[string]*fbreader.XMLTitleInfo)

	for i := 0; i < len(booksHash); i++ {
		go func(i int) {
			hash := booksHash[i]
			desc, err := stor.GetBookStorage().GetBookDesc(hash)
			if err == nil {
				srchStrBook := strings.ToLower(desc.Author.FirstName + desc.Author.LastName + desc.Author.MiddleName + desc.BookTitle + desc.Sequence.Name + desc.Keywords)
				for _, v := range fbreader.GetGenres(desc.Genre) {
					srchStrBook += strings.ToLower(v)
				}
				if strings.Contains(srchStrBook, srchStr) {
					booksmap[hash] = desc
					doneParse <- hash
				} else {
					doneParse <- "not found"
				}
			} else {
				doneParse <- "error"
			}
		}(i)
	}

	for i := 0; i < len(booksHash); i++ {
		<-doneParse
	}

	html := ""

	for hash, desc := range booksmap {

		html += `<a class="mainbookhref" rel="external" href="/book/` + hash + `">
				<div class="mainbookbox">
				<img class="mainbookcover" src="/getcover/` + hash + `" align="left" />
				<p>` + desc.BookTitle + `</p>
				<p>` + desc.Author.LastName + ` ` + desc.Author.MiddleName + ` ` + desc.Author.FirstName + `</p>
				<p>` + desc.Sequence.Name + ` ` + desc.Sequence.Number + `</p>
				<p><small>` + fbreader.ConvGenres(desc.Genre) + `</small></p>
				</div></a>`
	}
	res.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	return 200, html
}

func (l *Library) UserBooks(stor storage.Storages, req *http.Request, r render.Render, params martini.Params) {
	userName := params["user"]
	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login != 1 {
		log.Println("/login" + req.RequestURI)
		r.Redirect("/login" + req.RequestURI)
		return
	}

	user, err := stor.GetUserStorage().GetUser(userName)
	if err != nil {
		r.Error(404)
		log.Println("User not found", userName)
		return
	}

	books := getSortedBooksByDate(user.ReadingBooks)
	type MainPage struct {
		Books   []*fbreader.XMLTitleInfo
		Session *sessions.Session
		Edit    bool
	}
	edit := user.Name == ses.User.Name && ses.Login == 1 && len(books) > 0
	r.HTML(200, "library/userBooks", MainPage{books, ses, edit})
}

func (l *Library) UserRemoveBooks(stor storage.Storages, req *http.Request, res http.ResponseWriter) (int, string) {

	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login != 1 {
		return 401, "Зайдите на сайт под своим именем"
	}

	remBooks := []string{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&remBooks)
	if err != nil {
		return 400, ""
	}

	for _, b := range remBooks {
		ses.User.RemoveBookInfo(b)
	}
	stor.GetUserStorage().SaveUser(ses.User)
	res.Header().Set("Content-Type", "application/json")
	return 200, getMsg("")
}
