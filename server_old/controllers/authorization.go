package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"SimpleReader/server_old/models/captcha"
	"SimpleReader/server_old/models/crypt"
	"SimpleReader/server_old/models/sessions"
	"SimpleReader/server_old/models/storage"
	"SimpleReader/server_old/models/users"
	"SimpleReader/server_old/server"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type Authorization struct {
	server.Controller
}

func (h *Authorization) InitController(serv *server.Server) {
	serv.GetMartini().Router.Get("/login/?**", h.LoginPage)
	serv.GetMartini().Router.Post("/login", h.Verify)
	serv.GetMartini().Router.Post("/getkey", h.GetKey)
	serv.GetMartini().Router.Get("/logout/?**", h.Logout)
	serv.GetMartini().Router.Get("/registration", h.RegistrationPage)
	serv.GetMartini().Router.Post("/registration", h.Registration)
}

func GetBackTo(param string, req *http.Request) string {
	backto := param
	if backto == "" {
		backto = req.Referer()
		if len(backto) > 7 && backto[0:4] == "http" {
			backto = backto[strings.Index(backto[7:], "/")+7:]
		}
		if backto != "" {
			host := req.Header.Get("host")
			if host != "" {
				if !strings.Contains(strings.ToLower(backto), strings.ToLower(host)) {
					backto = "/"
				}
			}
		}
	}
	if backto == "" || backto == "/login" || backto == "/logout" {
		backto = "/"
	}
	if backto[0] != '/' {
		backto = "/" + backto
	}
	log.Println("back", backto, req.Referer())
	return backto
}

func (h *Authorization) LoginPage(stor storage.Storages, req *http.Request, r render.Render, param martini.Params) {
	backto := GetBackTo(param["_1"], req)
	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login != 1 {
		type LoginPage struct {
			BackTo  string
			Session *sessions.Session
		}
		r.HTML(200, "auth/login", LoginPage{backto, ses})
	} else {
		r.Redirect("/")
	}
}

func getVerifyMsg(code int, msg, captcha string) string {
	return `{
	"Code": ` + strconv.Itoa(code) + `,
	"Msg": "` + msg + `",
	"Captcha":"` + captcha + `"
	}`
}

func (h *Authorization) Verify(stor storage.Storages, req *http.Request) string {
	type userAuth struct {
		Login    string
		Password string
		Captcha  string
	}
	captchaImg := ""
	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login == 1 {
		getVerifyMsg(0, "", "")
	} else if ses.Login < 0 {
		ses.Login += 1
	}

	defer func() {
		if ses.Login == 0 {
			time.Sleep(3 * time.Second)
		}
	}()

	if ses.Key == nil {
		return getVerifyMsg(1, "Ошибка, повторите попытку позже", captchaImg)
	}

	auth := userAuth{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&auth)
	if err != nil {
		return getVerifyMsg(1, "Ошибка, неверные данные", captchaImg)
	}

	if ses.Login == 0 {
		if !captcha.VerifyCaptcha(ses.User, auth.Captcha) {
			captchaImg = captcha.GenerateCaptcha(ses.User)
			return getVerifyMsg(2, "Укажите правильное количество книг", captchaImg)
		}
	}

	if ses.Login == -1 || ses.Login == 0 {
		captchaImg = captcha.GenerateCaptcha(ses.User)
	}

	auth.Login = crypt.Decrypt(ses.Key, auth.Login)
	auth.Password = crypt.Decrypt(ses.Key, auth.Password)
	ses.Key = nil
	if auth.Login == "" || auth.Password == "" {
		return getVerifyMsg(2, "Имя пользователя или пароль не совпадают", captchaImg)
	}
	locUser, _ := stor.GetUserStorage().GetUser(auth.Login)
	if locUser == nil || locUser.Pass != auth.Password {
		return getVerifyMsg(2, "Имя пользователя и пароль не совпадают", captchaImg)
	}
	if locUser.Role == -2 {
		return getVerifyMsg(3, "Вы заблокированы", captchaImg)
	}

	for bookId, bookInfo := range ses.User.ReadingBooks {
		locBook, ok := locUser.ReadingBooks[bookId]
		if !ok || (bookInfo.LastRead.After(locBook.LastRead) && bookInfo.LastReadPage > locBook.LastReadPage) {
			locUser.ReadingBooks[bookId] = bookInfo
		}
	}
	ses.SetUser(locUser)
	ses.Login = 1

	return getVerifyMsg(0, "", "")
}

func (h *Authorization) GetKey(stor storage.Storages, req *http.Request, resp http.ResponseWriter) (int, string) {
	type pubkey struct {
		Mod string
		Exp string
	}

	ses := stor.GetSessionStorage().GetSession(req)
	ses.Key = crypt.GetKeyPair()
	jsKey := pubkey{}
	jsKey.Exp = fmt.Sprintf("%X", ses.Key.PublicKey.E)
	jsKey.Mod = fmt.Sprintf("%X", ses.Key.PublicKey.N)
	ret, err := json.Marshal(jsKey)
	if err == nil {
		resp.Header().Set("Content-Type", "application/json")
		return 200, string(ret)
	}
	return 500, ""
}

func (h *Authorization) Logout(stor storage.Storages, req *http.Request, r render.Render, param martini.Params) {
	backto := GetBackTo(param["_1"], req)
	ses := stor.GetSessionStorage().GetSession(req)
	ses.SetUser(users.NewUser())
	ses.Login = -5
	r.Redirect(backto)
}

func (h *Authorization) RegistrationPage(stor storage.Storages, req *http.Request, r render.Render, resp http.ResponseWriter) {

	ses := stor.GetSessionStorage().GetSession(req)
	if ses.Login == 1 {
		r.Redirect("/logout/registration")
		return
	}

	capt := captcha.GenerateCaptcha(ses.User)

	type Reg struct {
		Session *sessions.Session
		Captcha template.HTML
	}

	r.HTML(200, "auth/register", Reg{ses, template.HTML(capt)})
}

func (h *Authorization) Registration(stor storage.Storages, req *http.Request) string {
	type userReg struct {
		Login    string
		Password string
		Email    string
		Captcha  string
	}
	captchaImg := ""
	ses := stor.GetSessionStorage().GetSession(req)

	if ses.Key == nil {
		return getVerifyMsg(1, "Ошибка, повторите попытку позже", captchaImg)
	}

	reg := userReg{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reg)
	if err != nil {
		return getVerifyMsg(1, "Ошибка, неверные данные", captchaImg)
	}

	if !captcha.VerifyCaptcha(ses.User, reg.Captcha) || reg.Captcha == "" {
		captchaImg = captcha.GenerateCaptcha(ses.User)
		return getVerifyMsg(2, "Укажите правильное количество книг", captchaImg)
	}
	captchaImg = captcha.GenerateCaptcha(ses.User)

	reg.Login = crypt.Decrypt(ses.Key, reg.Login)
	reg.Password = crypt.Decrypt(ses.Key, reg.Password)
	reg.Email = crypt.Decrypt(ses.Key, reg.Email)
	ses.Key = nil
	if reg.Login == "" || reg.Password == "" {
		return getVerifyMsg(2, "Имя пользователя или пароль не могут быть пустыми", captchaImg)
	}

	user, _ := stor.GetUserStorage().GetUser(reg.Login)
	if user != nil {
		return getVerifyMsg(2, "Это имя занято, пожалуйста, выберите другое имя", captchaImg)
	}

	ok, _ := regexp.MatchString("[a-zA-Z0-9@#*_+]+", reg.Login)
	if len(reg.Login) > 30 || !ok {
		return getVerifyMsg(2, "Указанное имя недопустимо, содержит недопустимые символы или имеет длину более 30 знаков", captchaImg)
	}

	user = users.NewUser()
	user.Name = reg.Login
	user.Pass = reg.Password
	user.Email = reg.Email
	user.ReadingBooks = ses.User.ReadingBooks

	ses.SetUser(user)
	ses.Login = 1
	err = stor.GetUserStorage().SaveUser(user)

	if err != nil {
		log.Println("Error register user", reg, err)
		ses.Login = -5
		return getVerifyMsg(3, "Ошибка создания пользователя", captchaImg)
	}

	return getVerifyMsg(0, "", "")
}
