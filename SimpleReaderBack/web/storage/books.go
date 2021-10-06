package storage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/vdovindima/fb2"

	"SimpleReader/web/settings"
	"SimpleReader/web/utils"
)

// return hash and error
func AddBook(buf []byte) (string, error) {
	hash := utils.Md5HashBytes(buf)
	bookpath := path.Join(settings.Path, "books", hash)
	bookfile := path.Join(bookpath, "book.fb2")
	os.MkdirAll(bookpath, 0755)

	err := ioutil.WriteFile(bookfile, buf, 0664)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}

	// parse fb2
	book, err := fb2.New(buf).Unmarshal()
	if err != nil {
		os.RemoveAll(bookpath)
		return "", errors.New("Error parse book " + err.Error())
	}

	// save description
	desc, err := json.Marshal(book.Description)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", errors.New("Error parse book description")
	}

	err = ioutil.WriteFile(path.Join(bookpath, "desc.json"), desc, 0664)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}

	// save binaries
	for _, i := range book.Binary {
		buf, err := base64.StdEncoding.DecodeString(i.Value)
		if err == nil {
			ioutil.WriteFile(path.Join(bookpath, i.ID), buf, 0664)
		}
	}

	bins := book.Binary
	for i, _ := range bins {
		bins[i].Value = ""
	}

	bufBins, err := json.Marshal(bins)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", errors.New("Error parse book bins")
	}

	err = ioutil.WriteFile(path.Join(bookpath, "bins.json"), bufBins, 0664)
	if err != nil {
		os.RemoveAll(bookpath)
		return "", err
	}

	return hash, nil
}

func GetBook(hash string) (*fb2.FB2, error) {
	bookpath := path.Join(settings.Path, "books", hash)
	bookfile := path.Join(bookpath, "book.fb2")
	buf, err := ioutil.ReadFile(bookfile)
	if err != nil {
		return nil, err
	}
	book, err := fb2.New(buf).Unmarshal()
	if err != nil {
		return nil, errors.New("Error parse book " + err.Error())
	}
	return &book, nil
}

func GetDesc(hash string) (*fb2.Description, error) {
	infopath := path.Join(settings.Path, "books", hash, "desc.json")
	buf, err := ioutil.ReadFile(infopath)
	if err != nil {
		return nil, err
	}
	var desc *fb2.Description
	err = json.Unmarshal(buf, &desc)
	return desc, err
}

func GetBins(hash string) ([]fb2.Binary, error) {
	infopath := path.Join(settings.Path, "books", hash, "bins.json")
	buf, err := ioutil.ReadFile(infopath)
	if err != nil {
		return nil, err
	}
	var bins []fb2.Binary
	err = json.Unmarshal(buf, &bins)
	return bins, err
}

func GetImage(hash string, name string) ([]byte, error) {
	infopath := path.Join(settings.Path, "books", hash, name)
	return ioutil.ReadFile(infopath)
}
