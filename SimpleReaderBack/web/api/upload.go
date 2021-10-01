package api

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"

	"SimpleReader/web/models"
	"SimpleReader/web/session"
	"SimpleReader/web/storage"
)

func Upload(c *fiber.Ctx) error {
	ses := session.Get(c)
	user, _ := ses.Get("User").(*models.User)
	if user == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	defer form.RemoveAll()

	if files, ok := form.File["books"]; ok {
		loads := make([]string, 0)
		errloads := make([]string, 0)
		for _, f := range files {
			book, err := saveBook(f)
			if err != nil {
				errloads = append(errloads, f.Filename)
			} else {
				loads = append(loads, book)
			}
		}

		for _, b := range loads {
			if _, ok := user.ReadBooks[b]; !ok {
				user.ReadBooks[b] = models.BookInfo{BookHash: b}
			}
		}

		if len(errloads) > 0 {
			c.Status(200)
			c.JSON(errloads)
		}
	}
	return nil
}

func saveBook(file *multipart.FileHeader) (string, error) {
	Filename := file.Filename

	fileTmp, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileTmp.Close()
	buf, err := ioutil.ReadAll(fileTmp)
	if err != nil {
		return "", err
	}

	if strings.ToLower(path.Ext(Filename)) == ".zip" {
		buf, err = extractZip(bytes.NewReader(buf))
		if err != nil {
			return "", err
		}
	}
	if len(buf) == 0 {
		return "", errors.New("File is empty or error open file")
	}

	return storage.AddBook(buf)
}

func extractZip(breader *bytes.Reader) ([]byte, error) {
	zf, err := zip.NewReader(breader, int64(breader.Len()))
	if err != nil {
		return nil, err
	}
	for _, v := range zf.File {
		if strings.ToLower(path.Ext(v.Name)) == ".fb2" {
			r, err := v.Open()
			if err != nil {
				return nil, err
			}
			defer r.Close()
			return ioutil.ReadAll(r)
		}
	}
	return nil, errors.New("Not found fb2 in zip file")
}
