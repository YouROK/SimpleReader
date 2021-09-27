package api

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vdovindima/fb2"

	"SimpleReader/web/models"
	"SimpleReader/web/session"
	"SimpleReader/web/settings"
	"SimpleReader/web/utils"
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

	return saveFb2File(buf)
}

func saveFb2File(b []byte) (string, error) {
	hash := utils.Md5HashBytes(b)
	bookpath := path.Join(settings.Path, "books", hash)
	bookfile := path.Join(bookpath, "book.fb2")
	os.MkdirAll(bookpath, 0755)
	fb2File, err := os.Create(bookfile)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}
	fb2File.Write(b)
	fb2File.Close()

	book, err := fb2.New(b).Unmarshal()
	if err != nil {
		os.RemoveAll(bookpath)
		return "", errors.New("Error parse book")
	}

	// save description
	desc, err := json.Marshal(book.Description)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", errors.New("Error parse book description")
	}

	df, err := os.Create(path.Join(bookpath, "info"))
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}

	_, err = df.Write(desc)

	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}

	// save images
	for _, i := range book.Binary {
		imgF, err := os.Create(path.Join(bookpath, i.ID))
		if err == nil {
			buf, err := base64.StdEncoding.DecodeString(i.Value)
			if err == nil {
				imgF.Write(buf)
			}
			imgF.Close()
		}
	}
	return hash, nil
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
