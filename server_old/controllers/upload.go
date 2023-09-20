package controllers

import (
	"net/http"

	"github.com/martini-contrib/render"

	"SimpleReader/server_old/models/storage"
	"SimpleReader/server_old/models/users"
	"SimpleReader/server_old/server"
	// "strconv"
)

type Upload struct {
	server.Controller
}

func (upl *Upload) InitController(serv *server.Server) {
	serv.GetMartini().Router.Get("/upload", upl.UploadPage)
	serv.GetMartini().Router.Post("/upload", upl.UploadBook)
}

func (upl *Upload) UploadPage(stor storage.Storages, req *http.Request, r render.Render) {
	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login != 1 {
		r.Redirect("/login/upload")
		return
	}
	r.HTML(200, "upload/uploadPage", ses)
}

func (upl *Upload) UploadBook(stor storage.Storages, req *http.Request, res http.ResponseWriter, r render.Render) (int, string) {
	ses := stor.GetSessionStorage().GetSession(req)
	res.Header().Set("Content-Type", "text/plain")
	if ses.Login != 1 {
		return 401, "Ошибка загрузки, нужно войти на сайт под своим именем"
	}

	err := req.ParseMultipartForm(32 << 20) // 32 MB

	if err == nil {
		form := req.MultipartForm
		if form == nil || len(form.File) == 0 {
			return 400, "Выберите файлы для загрузки"
		} else {
			defer form.RemoveAll()
			if files, ok := form.File["fb2files"]; ok {
				loads := make([]string, 0)
				errloads := make([]string, 0)
				for _, f := range files {
					book, err := stor.GetBookStorage().WriteBook(f)
					if err != nil {
						errloads = append(errloads, f.Filename)
					} else {
						loads = append(loads, book)
					}
				}
				if len(errloads) > 0 {
					r.JSON(200, errloads)
				}
				for _, b := range loads {
					if _, ok := ses.User.ReadingBooks[b]; !ok {
						ses.User.ReadingBooks[b] = users.BookInfo{BookHash: b}
					}
				}
				stor.GetUserStorage().SaveUser(ses.User)
				return 200, ""
			}
			return 400, "Ошибка загрузки книги"
		}
	} else {
		return 500, "Ошибка загрузки, " + err.Error()
	}
}
