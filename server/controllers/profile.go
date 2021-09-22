package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/martini-contrib/render"

	"SimpleReader/server/models/storage"
	"SimpleReader/server/models/users"
	"SimpleReader/server/server"
)

type Styles struct {
	server.Controller
}

func (h *Styles) InitController(serv *server.Server) {
	serv.GetMartini().Router.Get("/styles", h.StylePage)
	serv.GetMartini().Router.Post("/setstyle", h.SetStyle)
}

func (h *Styles) StylePage(stor storage.Storages, req *http.Request, r render.Render) {
	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login != 1 {
		r.Redirect("/login/styles")
		return
	}
	r.HTML(200, "profile/stylePage", ses)
}

func (h *Styles) SetStyle(stor storage.Storages, req *http.Request, res http.ResponseWriter) (int, string) {
	ses := stor.GetSessionStorage().GetSession(req)
	res.Header().Set("Content-Type", "application/json")
	if ses.Login != 1 {
		return 200, getMsg("Вы не авторизованы")
	}

	style := users.Style{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&style)
	if err != nil {
		return 200, getMsg("Ошибка, неверные данные " + err.Error())
	}

	ses.User.Style = style
	stor.GetUserStorage().SaveUser(ses.User)

	return 200, getMsg("Настройки сохранены")
}

func getMsg(msg string) string {
	return `{"Message": "` + msg + `"}`
}
