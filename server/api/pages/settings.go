package pages

import (
	"SimpleReader/server/models/user"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetsPage(c *gin.Context) {
	ses := storage.GetSession(c)

	if !ses.IsLogin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "settings.gohtml", ses)
}

func SetStyle(c *gin.Context) {
	ses := storage.GetSession(c)

	if !ses.IsLogin() {
		c.JSON(http.StatusUnauthorized, getMsg("Вы не авторизованы"))
		return
	}

	style := user.Style{}
	err := c.BindJSON(&style)
	if err != nil {
		c.JSON(http.StatusBadRequest, getMsg("Ошибка, неверные данные "+err.Error()))
		return
	}

	ses.User.Style = &style
	storage.GetUserStorage().SaveUser(ses.User)

	c.JSON(200, getMsg("Сохранено"))
}

func getMsg(msg string) interface{} {
	return struct {
		Msg string
	}{msg}
}
