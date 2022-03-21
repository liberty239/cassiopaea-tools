package html

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func NewMetaNode(key, value string) *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		Data:     "meta",
		DataAtom: atom.Meta,
		Attr: []html.Attribute{
			{
				Key: "name",
				Val: key,
			},
			{
				Key: "content",
				Val: value,
			},
		},
	}
}

func NewMetaCharsetNode(value string) *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		Data:     "meta",
		DataAtom: atom.Meta,
		Attr: []html.Attribute{
			{
				Key: "charset",
				Val: value,
			},
		},
	}
}
