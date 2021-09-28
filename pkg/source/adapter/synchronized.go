package adapter

import (
	"io"
	"sync"
	"time"

	"github.com/liberty239/cassiopaea-tools/pkg/source"
)

type synchronized struct {
	source.Source
	mu sync.Mutex
}

func (s *synchronized) FetchAll(f func(time.Time, string, io.ReadCloser) error) error {
	return s.Source.FetchAll(func(ts time.Time, url string, rc io.ReadCloser) error {
		s.mu.Lock()
		defer s.mu.Unlock()
		return f(ts, url, rc)
	})
}

func Synchronized(source source.Source) source.Source {
	return &synchronized{Source: source}
}
