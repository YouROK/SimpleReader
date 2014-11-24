package fbreader

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"encoding/xml"
	"html/template"
	"log"
	"os"
	"utils"
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
	description *XMLTitleInfo
	xmlBook     *XMLBook
	xmlNotes    *XMLNotes
	xmlBinares  *[]XMLBinary
	annotation  *[]ContentItem
	bookcontent *[]ContentItem
	chapters    *[]Chapter
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

func (fbp *FBParser) GetDescription() *XMLTitleInfo {
	return fbp.description
}

func (fbp *FBParser) GetImages() *[]XMLBinary {
	return fbp.xmlBinares
}

func (fbp *FBParser) GetAnnotation() *[]ContentItem {
	if fbp.annotation == nil {
		return &[]ContentItem{}
	}
	return fbp.annotation
}

func (fbp *FBParser) GetContent() *[]ContentItem {
	return fbp.bookcontent
}

func (fbp *FBParser) GetContentCount() int {
	return len(*fbp.bookcontent)
}

func (fbp *FBParser) GetChapters() *[]Chapter {
	return fbp.chapters
}

func (fbp *FBParser) GetGenres() []string {
	return GetGenres(fbp.description.Genre)
}

func (fbp *FBParser) GetGenresStr() string {
	return ConvGenres(fbp.description.Genre)
}

func (fbp *FBParser) ParseFB() error {
	log.Println("Parse", fbp.bookpath)
	xmlFile, err := os.Open(fbp.bookpath)
	if err != nil {
		log.Println("Error open file", fbp.bookpath, err)
		return err
	}
	defer xmlFile.Close()
	fbp.xmlBook = &XMLBook{}
	fbp.xmlNotes = &XMLNotes{}
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(fbp.xmlBook)
	if err == nil {
		err = fbp.parseNotes()
	}
	if err == nil {
		fbp.xmlBinares = &fbp.xmlBook.Binarys
		for _, b := range *fbp.xmlBinares {
			b.Binary = ""
		}
		fbp.annotation = fbp.parseContent(fbp.xmlBook.Description.TitleInfo.Annotation.Content)
		fbp.description = &fbp.xmlBook.Description.TitleInfo
		fbp.chapters = &[]Chapter{}
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
