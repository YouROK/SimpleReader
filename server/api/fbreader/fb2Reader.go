package fbreader

import (
	"SimpleReader/server/api/utils"
	"SimpleReader/server/models/fb2"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"html/template"
	"log"
	"os"
)

type ContentItem struct {
	Text template.HTML
	Note template.HTML `json:"Note,omitempty"`
}

type Chapter struct {
	Name string
	Page int
}

type FBParser struct {
	bookpath    string
	bookhash    string
	description *fb2.XMLTitleInfo
	xmlBook     *fb2.XMLBook
	xmlNotes    *fb2.XMLNotes
	xmlBinares  []fb2.XMLBinary
	annotation  []ContentItem
	bookcontent []ContentItem
	chapters    []Chapter
}

func GetParser(bookpath string) *FBParser {
	fb := FBParser{}
	fb.bookpath = bookpath
	return &fb
}

func (fbp *FBParser) GetHash() string {
	if fbp.bookhash != "" {
		return fbp.bookhash
	}
	xmlFile, err := os.Open(fbp.bookpath)
	if err != nil {
		return ""
	}
	defer xmlFile.Close()
	fbp.bookhash = utils.Md5HashReader(xmlFile)
	return fbp.bookhash
}

func (fbp *FBParser) GetPath() string {
	return "/books/" + fbp.GetHash()
}

func (fbp *FBParser) GetDescription() *fb2.XMLTitleInfo {
	return fbp.description
}

func (fbp *FBParser) GetImages() []fb2.XMLBinary {
	return fbp.xmlBinares
}

func (fbp *FBParser) GetAnnotation() []ContentItem {
	if fbp.annotation == nil {
		return []ContentItem{}
	}
	return fbp.annotation
}

func (fbp *FBParser) GetContent() []ContentItem {
	return fbp.bookcontent
}

func (fbp *FBParser) GetContentCount() int {
	return len(fbp.bookcontent)
}

func (fbp *FBParser) GetChapters() []Chapter {
	return fbp.chapters
}

func (fbp *FBParser) GetGenres() []string {
	return utils.GetGenres(fbp.description.Genre)
}

func (fbp *FBParser) GetGenresLine() string {
	return utils.LineGenres(fbp.description.Genre)
}

func (fbp *FBParser) ParseFB() error {
	xmlFile, err := os.Open(fbp.bookpath)
	if err != nil {
		log.Println("Error open file", fbp.bookpath, err)
		return err
	}
	defer xmlFile.Close()
	fbp.xmlBook = &fb2.XMLBook{}
	fbp.xmlNotes = &fb2.XMLNotes{}
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(fbp.xmlBook)
	if err == nil {
		err = fbp.parseNotes()
	}
	if err == nil {
		fbp.xmlBinares = fbp.xmlBook.Binarys
		for _, b := range fbp.xmlBinares {
			b.Binary = ""
		}
		fbp.annotation = fbp.parseContent(fbp.xmlBook.Description.TitleInfo.Annotation.Content)
		fbp.description = &fbp.xmlBook.Description.TitleInfo
		fbp.chapters = []Chapter{}
		for _, b := range fbp.xmlBook.Bodys {
			if b.Type != "notes" {
				fbp.bookcontent = fbp.parseContent(b.Content)
			}
		}
	}
	fbp.xmlBook = nil
	fbp.xmlNotes = nil
	if err != nil {
		log.Println("Error parse", err)
	}
	return err
}
