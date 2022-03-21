package html

import (
	"bytes"
	"io"
	"strings"

	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
	"golang.org/x/text/unicode/norm"
)

func RenderNode(w io.Writer, node *html.Node, pretty bool) error {
	return RenderNodes(w, []*html.Node{node}, pretty)
}

func sanitizeNode(n *html.Node) *html.Node {
	switch n.Type {
	case html.TextNode, html.CommentNode:
		n.Data = strings.ReplaceAll(n.Data, "“", "\"")
		n.Data = strings.ReplaceAll(n.Data, "”", "\"")
		n.Data = strings.ReplaceAll(n.Data, "’", "'")
		n.Data = strings.Trim(n.Data, "â€œ�\u009d")
		n.Data = norm.NFC.String(n.Data)
	}
	return n
}

func RenderNodes(w io.Writer, nodes []*html.Node, pretty bool) error {
	var b bytes.Buffer

	if len(nodes) > 0 {
		for c := nodes[0].FirstChild; c != nil; c = c.NextSibling {
			err := html.Render(&b, sanitizeNode(c))
			if err != nil {
				return err
			}
		}
	}

	r := strings.NewReader(gohtml.Format(b.String()))
	_, err := io.Copy(w, r)
	return err
}
