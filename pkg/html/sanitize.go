package html

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/go-shiori/go-readability"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

func SanitizeNodes(nodes []*html.Node) (*html.Node, error) {
	var b bytes.Buffer

	if err := RenderNodes(&b, nodes, false); err != nil {
		return nil, err
	}

	art, err := readability.FromReader(&b, nil)
	if err != nil {
		return nil, err
	}

	// Replace line breaks with paragraphs.
	art.Content = strings.ReplaceAll(art.Content, "<br/>", "</p><p>")
	// Replace unneded whitespaces.
	art.Content = regexp.MustCompile(`\s+`).ReplaceAllString(art.Content, " ")

	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowImages()
	p.AllowLists()
	p.AllowTables()
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowAttrs("cite").OnElements("blockquote")
	p.AllowElements("br", "div", "hr", "p")
	p.AllowAttrs("href").OnElements("a")
	p.AllowAttrs("title")
	return html.Parse(bytes.NewBufferString(p.Sanitize(art.Content)))
}
