package fbreader

import (
	"encoding/xml"
	"html"
	"html/template"
	"strings"
)

func (fbp *FBParser) parseNotes() error {
	for _, n := range fbp.xmlBook.Bodys {
		if n.Type == "notes" {
			decoder := xml.NewDecoder(strings.NewReader("<root>" + n.Content + "</root>"))
			return decoder.Decode(fbp.xmlNotes)
		}
	}
	return nil
}

func (fbp *FBParser) parseContent(xmlText string) *[]ContentItem {
	if xmlText == "" {
		return nil
	}
	root := Parse(xmlText)
	var xmlList []Node = []Node{}
	fbp.getNodeList(root, &xmlList)
	return fbp.parseNodeList(xmlList)
}

func (fbp *FBParser) getNodeList(node *Node, ret *[]Node) {
	for _, n := range node.Childs {
		switch strings.ToLower(n.Key) {
		case "style", "author", "text-author", "date", "empty-line", "stanza", "cite", "poem", "image", "title", "epigraph", "subtitle", "emphasis", "p":
			{
				*ret = append(*ret, n)
			}
		case "section":
			{
				sectionTitle := ""
				for _, t := range n.Childs {
					if len(t.Childs) > 0 && len(t.Childs[0].Childs) > 0 {
						if t.Key == "title" && t.Childs[0].Childs[0].Key == "#text" {
							sectionTitle = t.Childs[0].Childs[0].Value
						}
					}
				}
				if sectionTitle != "" {
					*ret = append(*ret, Node{Key: "section", Value: sectionTitle, Parent: n.Parent})
				}
				fbp.getNodeList(&n, ret)
			}
		default:
			{
				fbp.getNodeList(&n, ret)
			}
		}

	}
}

func (fbp *FBParser) parseNodeList(list []Node) *[]ContentItem {
	if len(list) == 0 {
		return nil
	}
	ret := []ContentItem{}
	for _, v := range list {
		var htm string
		var notes string
		fbp.xml2html([]Node{v}, &htm, &notes)
		ret = append(ret, ContentItem{Text: template.HTML(htm), Note: template.HTML(notes)})
		if v.Key == "section" {
			*fbp.chapters = append(*fbp.chapters, Chapter{v.Value, len(ret) - 1})
		}
	}
	return &ret
}

func (fbp *FBParser) xml2html(node []Node, ret *string, notes *string) {
	for _, v := range node {
		switch strings.ToLower(v.Key) {
		case "date", "style", "emphasis":
			{
				*ret += "<i>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</i>\n"
			}
		case "p", "subtitle":
			{
				*ret += "<p>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</p>\n"
			}
		case "strong":
			{
				*ret += "<b>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</b>\n"
			}
		case "empty-line":
			{
				*ret += "<br/>"
			}
		case "epigraph", "stanza", "poem":
			{
				*ret += "<div align=\"center\">\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</div>\n"
			}
		case "cite":
			{
				*ret += "<cite>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</cite>\n"
			}
		case "text-author":
			{
				*ret += "<b><p>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</p></b>\n"
			}
		case "book-title":
			{
				*ret += "<h2><p>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</p></h2>\n"
			}
		case "#text":
			{
				txtarr := strings.Split(v.Value, "\n")
				txt := ""
				for _, t := range txtarr {
					txt += strings.TrimSpace(t) + " "
				}
				*ret += html.EscapeString(strings.TrimSpace(txt))
			}
		case "author":
			{
				*ret += "<h3><p>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</p></h3>\n"
			}
		case "title":
			{
				*ret += "<h2>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</h2>\n"
			}
		case "v":
			{
				*ret += "<div class=\"paragraph\">\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</div>\n"
			}
		case "section":
			{
				*ret += "<div id=\"" + v.Value + "\"></div>\n"
			}
		case "code":
			{
				*ret += "<code>\n"
				fbp.xml2html(v.Childs, ret, notes)
				*ret += "</code>\n"
			}

		case "a":
			{
				noteLink := fbp.getAttr(v.Attr, "href")
				noteType := fbp.getAttr(v.Attr, "type")
				if noteLink != "" && noteType != "" {
					if noteType == "note" {
						noteid := noteLink[1:]
						for _, n := range fbp.xmlNotes.Notes {
							if n.Id == noteid {
								*notes += "<div><span><i><p>" + n.Title + "</p></i></span><p>" + n.Note + "</p></div><br/>\n"
								*ret += "<abbr title=\"" + n.Note + "\"><sup>"
								fbp.xml2html(v.Childs, ret, notes)
								*ret += "</sup></abbr>\n"
							}
						}

					} else {
						fbp.xml2html(v.Childs, ret, notes)
					}
				}
			}
		case "image":
			{
				imgLink := fbp.getAttr(v.Attr, "href")
				if imgLink != "" {
					*ret += "<div align=\"center\"><img class=\"pageimage\" src=\"/getimage/" + fbp.GetHash() + "/" + imgLink[1:] + "\" /></div><br/>"
				}
			}
		default:
			{
				fbp.xml2html(v.Childs, ret, notes)
				*ret += " "
			}
		}
	}
}

func (fbp *FBParser) getAttr(attr []xml.Attr, name string) string {
	for _, v := range attr {
		if v.Name.Local == name {
			return v.Value
		}
	}
	return ""
}
