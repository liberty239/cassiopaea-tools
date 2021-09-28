package adapter

import (
	"io"
	"sort"
	"time"

	"github.com/liberty239/cassiopaea-tools/pkg/source"
)

type ascending struct {
	source.Source
}

func (a *ascending) FetchAll(f func(time.Time, string, io.ReadCloser) error) error {
	type data struct {
		ts  time.Time
		uri string
		rc  io.ReadCloser
	}

	var xs []data
	err := a.Source.
		FetchAll(func(ts time.Time, uri string, rc io.ReadCloser) error {
			xs = append(xs, data{ts: ts, uri: uri, rc: rc})
			return nil
		})
	if err != nil {
		return err
	}

	sort.Slice(xs, func(i, j int) bool {
		return xs[i].ts.Before(xs[j].ts)
	})

	for _, x := range xs {
		if err := f(x.ts, x.uri, x.rc); err != nil {
			return err
		}
	}
	return nil
}

func Ascending(source source.Source) source.Source {
	return &ascending{Source: source}
}
