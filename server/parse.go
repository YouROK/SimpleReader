package server

import (
	"SimpleReader/server/api/fbreader"
	"SimpleReader/server/api/utils"
	"SimpleReader/server/models/global"
	"SimpleReader/server/storage"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func ParseDir() {
	dir := filepath.Join(storage.GetStorage().GetPath(), "parse")

	for !global.Stoped {
		doneParse := make(chan string)
		go func() {
			filepath.Walk(dir,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					ext := strings.ToLower(filepath.Ext(path))
					if ext == ".fb2" || ext == ".zip" {
						err = addBookFile(path)
						if err != nil {
							os.Remove(path)
							log.Println("Error add book:", err.Error())
						} else {
							os.Remove(path)
							doneParse <- path
						}
					}
					return nil
				})
			close(doneParse)
		}()

		for range doneParse {
		}

		time.Sleep(time.Minute * 1)
	}
}

func addBookFile(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var bufs [][]byte
	if strings.ToLower(filepath.Ext(path)) == ".zip" {
		bufs, err = storage.ExtractZip(bytes.NewReader(buf))
		if err != nil {
			return err
		}
	} else {
		bufs = [][]byte{buf}
	}
	if len(bufs) == 0 {
		return errors.New("File is empty or error open file")
	}

	for _, b := range bufs {
		_, e := saveFb2File(b)
		if e != nil {
			log.Println("Error save book:", e)
			err = e
		}
	}
	return err
}

func saveFb2File(b []byte) (string, error) {
	hash := utils.Md5HashBytes(b)
	bookpath := path.Join(storage.GetStorage().GetPath(), "books", hash)
	bookfile := path.Join(bookpath, "book.fb2")
	os.MkdirAll(bookpath, 0755)
	fb2File, err := os.Create(bookfile)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}
	fb2File.Write(b)
	fb2File.Close()

	fb2 := fbreader.GetParser(bookfile)
	err = fb2.ParseFB()
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}

	//save description
	desc, err := json.Marshal(fb2.GetDescription())
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

	//save img
	bins := fb2.GetImages()
	for _, i := range bins {
		imgF, err := os.Create(path.Join(bookpath, i.Id))
		if err == nil {
			buf, err := base64.StdEncoding.DecodeString(i.Binary)
			if err == nil {
				imgF.Write(buf)
			}
			imgF.Close()
		}
	}
	return hash, nil
}
