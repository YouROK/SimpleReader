package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vdovindima/fb2"

	"SimpleReader/web/models"
	"SimpleReader/web/session"
	"SimpleReader/web/settings"
	"SimpleReader/web/storage"
)

func Search(c *fiber.Ctx) error {
	ses := session.Get(c)
	user, _ := ses.Get("User").(*models.User)
	if user == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	query := c.Query("query")
	query = strings.ToLower(query)

	books := getAllBooksDesc()
	var matches []string
	for hash, book := range books {
		buf, _ := json.Marshal(book)
		if strings.Index(strings.ToLower(string(buf)), query) != -1 {
			matches = append(matches, hash)
		}
	}
	return c.JSON(matches)
}

func getAllBooksDesc() map[string]*fb2.Description {
	bookspath := path.Join(settings.Path, "books")
	dirs, err := ioutil.ReadDir(bookspath)
	if err != nil {
		return nil
	}
	list := map[string]*fb2.Description{}
	for _, dir := range dirs {
		desc, err := storage.GetDesc(dir.Name())
		if err == nil {
			list[dir.Name()] = desc
		}
	}
	return list
}
