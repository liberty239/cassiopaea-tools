package forum

import (
	"bytes"
	"context"
	"io"
	"sort"
	"strings"
	"time"

	xhtml "golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"

	"github.com/liberty239/cassiopaea-tools/pkg/html"
	"github.com/liberty239/cassiopaea-tools/pkg/http"
	"github.com/liberty239/cassiopaea-tools/pkg/source"
)

type forum struct{}

func (*forum) List() (ret []time.Time, err error) {
	s := scraper{
		onSessionLink: func(e *colly.HTMLElement) error {
			if p := timestampFromPath(e.Attr("href")); len(p) > 0 {
				var ts time.Time
				ts, err = parseTimestamp(p)
				if err != nil {
					return err
				}
				ret = append(ret, ts)
			}
			return nil
		},
	}
	if err = s.do(allSessionsSearchQueryUrls()...); err != nil {
		return
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Before(ret[j])
	})
	return
}

func is22Feb2010(e *colly.HTMLElement) bool {
	return strings.HasPrefix(
		e.Request.URL.Path,
		"/forum/threads/session-22-february-2010",
	)
}

type readCloser struct {
	io.Reader
}

func (rc *readCloser) Close() error {
	return nil
}

func elementToDocument(ctx context.Context, e *colly.HTMLElement) (ret io.ReadCloser, ts time.Time, err error) {
	if is22Feb2010(e) {
		if e.Index != 1 {
			return
		}
	} else {
		if e.Index != 0 {
			return
		}
	}

	if len(e.DOM.Nodes) > 0 {
		if p := timestampFromPath(e.Request.URL.Path); len(p) > 0 {
			ts, err = parseTimestamp(p)
			if err != nil {
				return
			}
		}

		var doc *xhtml.Node
		doc, err = html.SanitizeNodes(e.DOM.Nodes)
		if err != nil {
			return
		}

		q := goquery.NewDocumentFromNode(doc)
		q.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if src, ok := s.Attr("src"); ok {
				ctx, cf := context.WithTimeout(ctx, scrapeTimeout)
				defer cf()

				src, err = http.GetDataURL(ctx, absoluteUrl(src))
				if err != nil {
					return false
				}

				s.SetAttr("src", src)
			}
			return true
		})
		if err != nil {
			return
		}

		if n := q.Find("head").Nodes; len(n) > 0 {
			for _, meta := range []*xhtml.Node{
				html.NewMetaNode("url", e.Request.URL.String()),
				html.NewMetaNode("timestamp", ts.Format(time.RFC3339)),
			} {
				n[0].AppendChild(meta)
			}
		}

		var b bytes.Buffer
		ret = &readCloser{Reader: &b}
		err = html.RenderNodes(&b, q.Nodes, true)
		return
	}
	return
}

func (*forum) Fetch(ts time.Time) (r io.ReadCloser, uri string, err error) {
	s := scraper{
		onSession: func(e *colly.HTMLElement) error {
			if r == nil {
				r, _, err = elementToDocument(context.TODO(), e)
				uri = e.Request.URL.String()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	err = s.do(
		newSearchQuery(ts.Format("Session 2 January 2006")),
		newSearchQuery(ts.Format("Sesssion 2 January 2006")),
	)
	return
}

func (*forum) FetchAll(f func(time.Time, string, io.ReadCloser) error) error {
	s := scraper{
		onSession: func(e *colly.HTMLElement) error {
			r, ts, err := elementToDocument(context.TODO(), e)
			if err != nil {
				return err
			}

			if r != nil {
				return f(ts, e.Request.URL.String(), r)
			}
			return nil
		},
	}
	if err := s.do(allSessionsSearchQueryUrls()...); err != nil {
		return err
	}

	return nil
}

func NewSource() source.Source {
	return new(forum)
}
