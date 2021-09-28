package source

import (
	"io"
	"time"
)

type Source interface {
	List() ([]time.Time, error)
	Fetch(ts time.Time) (io.ReadCloser, string, error)
	FetchAll(f func(time.Time, string, io.ReadCloser) error) error
}
