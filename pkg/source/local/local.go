package local

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/liberty239/cassiopaea-tools/pkg/source"
)

type local struct {
	path string
}

type entry struct {
	timestamp time.Time
	url       string
	path      string
}

func metaTag(doc *goquery.Document, name string) (ret string) {
	doc.
		Find(fmt.Sprintf("meta[name=%s]", name)).
		EachWithBreak(func(_ int, s *goquery.Selection) bool {
			var ok bool
			ret, ok = s.Attr("content")
			return ok
		})
	return
}

func readMetaTags(doc *goquery.Document) (entry, error) {
	ts, err := time.Parse(time.RFC3339, metaTag(doc, "timestamp"))
	if err != nil {
		return entry{}, err
	}
	return entry{
		timestamp: ts,
		url:       metaTag(doc, "url"),
	}, nil
}

func (l *local) readFiles() (ret []entry, err error) {
	err = filepath.Walk(l.path, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		switch strings.ToLower(filepath.Ext(info.Name())) {
		case ".htm", ".html":
		default:
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(f)
		if err != nil {
			f.Close()
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		e, err := readMetaTags(doc)
		if err != nil {
			return nil
		}

		e.path = path

		ret = append(ret, e)
		return nil
	})
	return
}

func (l *local) List() (ret []time.Time, err error) {
	xs, err := l.readFiles()
	if err != nil {
		return nil, err
	}

	for _, x := range xs {
		ret = append(ret, x.timestamp)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Before(ret[j])
	})
	return
}

func (l *local) Fetch(ts time.Time) (r io.ReadCloser, uri string, err error) {
	xs, err := l.readFiles()
	if err != nil {
		return nil, "", err
	}

	for _, x := range xs {
		if x.timestamp.Equal(ts) {
			f, err := os.Open(x.path)
			return f, x.url, err
		}
	}

	return nil, "", errors.New("not found")
}

type lazyFile struct {
	path string
	f    *os.File
}

func (l *lazyFile) Read(p []byte) (int, error) {
	if l.f == nil {
		f, err := os.Open(l.path)
		if err != nil {
			return 0, err
		}
		l.f = f
	}
	return l.f.Read(p)
}

func (l *lazyFile) Close() error {
	if l.f == nil {
		return nil
	}
	return l.f.Close()
}

func (l *local) FetchAll(f func(time.Time, string, io.ReadCloser) error) error {
	xs, err := l.readFiles()
	if err != nil {
		return err
	}

	for _, x := range xs {
		if err := f(x.timestamp, x.url, &lazyFile{path: x.path}); err != nil {
			return err
		}
	}
	return nil
}

func NewSource(path string) source.Source {
	return &local{
		path: path,
	}
}
