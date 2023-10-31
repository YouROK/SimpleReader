package fbreader

import (
	"encoding/xml"
	"strings"
)

type Node struct {
	Key    string
	Value  string
	Attr   []xml.Attr
	Childs []*Node
	Parent *Node
}

func Parse(xmlStr string) *Node {
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	root := Node{}
	root.Childs = make([]*Node, 0)
	parse(decoder, &root)
	return &root
}

func parse(decoder *xml.Decoder, root *Node) {
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch elm := token.(type) {
		case xml.StartElement:
			{
				next := Node{}
				next.Parent = root
				next.Key = elm.Name.Local
				if len(elm.Attr) > 0 {
					next.Attr = elm.Copy().Attr
				}
				parse(decoder, &next)
				root.Childs = append(root.Childs, &next)
			}
		case xml.EndElement:
			{
				return
			}
		case xml.CharData:
			{
				txt := strings.TrimSpace(string([]byte(elm)))
				if len(txt) > 0 {
					next := Node{}
					next.Parent = root
					next.Key = "#text"
					next.Value = txt
					root.Childs = append(root.Childs, &next)
				}
			}
		}
	}
}
