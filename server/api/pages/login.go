package pages

import (
	"SimpleReader/server/api/captcha"
	"SimpleReader/server/api/crypt"
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func getBackTo(c *gin.Context) string {
	backto := c.GetHeader("Referer")
	if len(backto) > 7 && backto[0:4] == "http" {
		backto = backto[strings.Index(backto[7:], "/")+7:]
	}
	if backto != "" {
		host := c.GetHeader("host")
		if host != "" {
			if !strings.Contains(strings.ToLower(backto), strings.ToLower(host)) {
				backto = "/"
			}
		}
	}
	if backto == "" || backto == "/login" || backto == "/logout" {
		backto = "/"
	}
	if backto[0] != '/' {
		backto = "/" + backto
	}
	log.Println("back", backto, c.GetHeader("Referer"))
	return backto
}

func LoginPage(c *gin.Context) {
	backPage := getBackTo(c)
	c.HTML(http.StatusOK, "login.gohtml", backPage)
}

type respLogin struct {
	Code    int
	Msg     string
	Captcha string
}

func newRespLogin(code int, msg, cap string) respLogin {
	return respLogin{
		Code:    code,
		Msg:     msg,
		Captcha: cap,
	}
}

type reqLogin struct {
	Login    string
	Password string
	Captcha  string
}

func Login(c *gin.Context) {
	ses := storage.GetSession(c)
	if ses.Login < 0 {
		ses.Login += 1
	}

	defer func() {
		if ses.Login == 0 {
			time.Sleep(3 * time.Second)
		}
	}()

	if ses.Key == nil {
		c.JSON(200, newRespLogin(1, "Ошибка, повторите попытку позже", ""))
		return
	}

	auth := reqLogin{}
	err := c.BindJSON(&auth)
	if err != nil {
		c.JSON(200, newRespLogin(1, "Ошибка, неверные данные", ""))
		return
	}

	captchaImg := auth.Captcha

	if ses.Login == 0 {
		if !captcha.VerifyCaptcha(ses.User, auth.Captcha) {
			captchaImg = captcha.GenerateCaptcha(ses.User)
			c.JSON(200, newRespLogin(2, "Укажите правильное число на капче", captchaImg))
			return
		}
	}

	if ses.Login == -1 || ses.Login == 0 {
		captchaImg = captcha.GenerateCaptcha(ses.User)
	}

	auth.Login = crypt.Decrypt(ses.Key, auth.Login)
	auth.Password = crypt.Decrypt(ses.Key, auth.Password)
	ses.Key = nil
	if auth.Login == "" || auth.Password == "" {
		c.JSON(200, newRespLogin(2, "Имя пользователя или пароль не совпадают", captchaImg))
		return
	}
	locUser, _ := storage.GetUserStorage().GetUser(auth.Login)
	if locUser == nil || locUser.Pass != auth.Password {
		c.JSON(200, newRespLogin(2, "Имя пользователя и пароль не совпадают", captchaImg))
		return
	}
	if locUser.Role == -2 {
		c.JSON(200, newRespLogin(3, "Вы заблокированы", captchaImg))
		return
	}

	ses.SetUser(locUser)
	ses.Login = 1

	refer := c.GetHeader("Referer")
	if refer == "" || filepath.Base(refer) == "login" {
		refer = "/"
	}
	c.JSON(200, newRespLogin(0, refer, ""))
}

func GetKey(c *gin.Context) {
	type pubkey struct {
		Mod string
		Exp string
	}

	ses := storage.GetSessionStorage().GetSession(c)
	ses.Key = crypt.GetKeyPair()
	jsKey := pubkey{}
	jsKey.Exp = fmt.Sprintf("%X", ses.Key.PublicKey.E)
	jsKey.Mod = fmt.Sprintf("%X", ses.Key.PublicKey.N)
	c.JSON(200, jsKey)
}

func Logout(c *gin.Context) {
	backPage := getBackTo(c)
	ses := storage.GetSessionStorage().GetSession(c)
	ses.SetUser(user.NewUser())
	ses.Login = -5
	c.Redirect(http.StatusFound, backPage)
}
