package storage

import (
	"SimpleReader/server/api/fbreader"
	"SimpleReader/server/api/utils"
	"SimpleReader/server/models/fb2"
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type BookStorage struct {
	path      string
	booksList map[string]*fbreader.FBParser
}

func newBookStorage(path string) *BookStorage {
	b := BookStorage{}
	b.path = path
	b.booksList = make(map[string]*fbreader.FBParser)
	return &b
}

func (bs *BookStorage) Exit() {

}

func (bs *BookStorage) GetPath() string {
	return bs.path
}

func (bs *BookStorage) GetBooks() []string {
	dirs, err := os.ReadDir(bs.path)
	var ret []string
	if err == nil {
		for _, d := range dirs {
			if d.IsDir() && d.Name() != ".gitignore" {
				ret = append(ret, d.Name())
			}
		}
		return ret
	}
	return nil
}

func (bs *BookStorage) GetBook(hash string) *fbreader.FBParser {
	if b, ok := bs.booksList[hash]; ok {
		return b
	}
	fb := fbreader.GetParser(path.Join(bs.path, hash, "book.fb2"))
	err := fb.ParseFB()
	if err != nil {
		return nil
	}
	bs.booksList[hash] = fb
	return fb
}

func (bs *BookStorage) GetBookDesc(hash string) (*fb2.XMLTitleInfo, error) {
	if !bs.BookExists(hash) {
		return nil, errors.New("Book not found")
	}
	bookfile := path.Join(bs.GetPath(), hash, "info")
	file, err := os.Open(bookfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	desc := fb2.XMLTitleInfo{}
	js := json.NewDecoder(file)
	err = js.Decode(&desc)
	desc.Hash = hash
	return &desc, err
}

func (bs *BookStorage) GetBookImgs(hash string) ([]string, error) {
	if !bs.BookExists(hash) {
		return nil, errors.New("Book not found")
	}
	bookpath := path.Join(bs.GetPath(), hash)
	files, err := ioutil.ReadDir(bookpath)
	if err != nil {
		return nil, err
	}
	imgs := make([]string, 0)
	for _, v := range files {
		if v.Name() != "info" && v.Name() != "book.fb2" {
			imgs = append(imgs, v.Name())
		}
	}
	return imgs, nil
}

func (bs *BookStorage) BookExists(hash string) bool {
	_, err := os.Stat(path.Join(bs.path, hash))
	return nil == err
}

//save book

func ExtractZip(breader *bytes.Reader) ([][]byte, error) {
	zf, err := zip.NewReader(breader, int64(breader.Len()))
	if err != nil {
		return nil, err
	}
	var books [][]byte
	for _, v := range zf.File {
		if strings.ToLower(path.Ext(v.Name)) == ".fb2" {
			r, err := v.Open()
			if err != nil {
				return nil, err
			}
			buf, err := io.ReadAll(r)
			r.Close()
			if err != nil {
				return nil, err
			}
			books = append(books, buf)
		}
	}
	if len(books) == 0 {
		return nil, errors.New("Not found fb2 in zip file")
	}
	return books, err
}

func (bs *BookStorage) WriteBook(file *multipart.FileHeader) ([]string, error) {
	Filename := file.Filename

	fileTmp, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileTmp.Close()
	buf, err := io.ReadAll(fileTmp)
	if err != nil {
		return nil, err
	}

	var bufs [][]byte

	if strings.ToLower(filepath.Ext(Filename)) == ".zip" {
		bufs, err = ExtractZip(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
	} else {
		bufs = [][]byte{buf}
	}
	if len(bufs) == 0 {
		return nil, errors.New("File is empty or error open file")
	}

	var hashs []string
	for _, b := range bufs {
		hash, e := bs.saveFb2File(b)
		if e != nil {
			log.Println("Error save book:", e)
			err = e
		} else {
			hashs = append(hashs, hash)
		}
	}
	return hashs, err
}

func (bs *BookStorage) saveFb2File(b []byte) (string, error) {
	hash := utils.Md5HashBytes(b)
	bookpath := path.Join(bs.GetPath(), hash)
	bookfile := path.Join(bookpath, "book.fb2")
	os.MkdirAll(bookpath, 0755)
	fb2File, err := os.Create(bookfile)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}
	fb2File.Write(b)
	fb2File.Close()

	fb2 := bs.GetBook(hash)

	if fb2 == nil {
		os.RemoveAll(bookpath)
		return "", errors.New("Error parse book")
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
