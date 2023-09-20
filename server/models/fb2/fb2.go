package fb2

import (
	"SimpleReader/server/api/utils"
	"encoding/xml"
)

type XMLAuthor struct {
	FirstName  string `xml:"first-name"`
	MiddleName string `xml:"middle-name"`
	LastName   string `xml:"last-name"`
}

type XMLSequence struct {
	Name   string `xml:"name,attr"`
	Number string `xml:"number,attr"`
}

type XMLText struct {
	Content string `xml:",innerxml"`
}

type XMLImage struct {
	ImgLink string `xml:"href,attr"`
}

type XMLTitleInfo struct {
	XMLName    xml.Name    `xml:"title-info" json:"-"`
	Hash       string      `xml:"-" json:"-"`
	Genre      []string    `xml:"genre"`
	Author     XMLAuthor   `xml:"author"`
	BookTitle  string      `xml:"book-title"`
	Annotation XMLText     `xml:"annotation"`
	Keywords   string      `xml:"keywords"`
	Date       string      `xml:"date"`
	Coverpage  XMLImage    `xml:"coverpage>image"`
	Lang       string      `xml:"lang"`
	Translator XMLAuthor   `xml:"translator"`
	Sequence   XMLSequence `xml:"sequence"`
}

func (ti *XMLTitleInfo) GetGenresStr() string {
	return utils.LineGenres(ti.Genre)
}

type XMLDescription struct {
	TitleInfo XMLTitleInfo
}

type XMLBody struct {
	Type    string `xml:"name,attr"`
	Content string `xml:",innerxml"`
}

type XMLBinary struct {
	ContentType string `xml:"content-type,attr"`
	Id          string `xml:"id,attr"`
	Binary      string `xml:",chardata"`
}

type XMLBook struct {
	XMLName     xml.Name       `xml:"FictionBook"`
	Description XMLDescription `xml:"description"`
	Bodys       []XMLBody      `xml:"body"`
	Binarys     []XMLBinary    `xml:"binary"`
}

type XMLNote struct {
	Id    string `xml:"id,attr"`
	Title string `xml:"title>p"`
	Note  string `xml:"p"`
}

type XMLNotes struct {
	Notes []XMLNote `xml:"section"`
}
