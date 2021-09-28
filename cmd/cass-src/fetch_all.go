package main

import (
	"io"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/pterm/pterm"

	"github.com/liberty239/cassiopaea-tools/pkg/term"
)

func doFetchAll(args args) error {
	then := time.Now()

	s, err := term.Spinner().Start("Fetching...")
	if err != nil {
		s.Fail(err)
		return err
	}

	if err := os.MkdirAll(args.Directory, 0755); err != nil {
		s.Fail(err)
		return err
	}

	var count int32
	err = argsToSource(args).
		FetchAll(func(ts time.Time, _ string, r io.ReadCloser) error {
			err := func() error {
				defer r.Close()

				f, err := os.Create(
					path.Join(
						args.Directory,
						ts.Format("2006-01-02.html"),
					),
				)
				if err != nil {
					return err
				}

				if _, err = io.Copy(f, r); err != nil {
					f.Close()
					return err
				}

				return f.Close()
			}()
			if err != nil {
				pterm.Error.Println(ts.Format("02 January 2006"))
			} else {
				pterm.Success.Println(ts.Format("02 January 2006"))
				atomic.AddInt32(&count, 1)
			}
			return err
		})
	if err != nil {
		s.Fail(err)
		return err
	}
	s.Success("Sessions fetched: ", atomic.LoadInt32(&count))
	pterm.Info.Println("Fetch took", time.Since(then).Seconds(), "seconds")
	return nil
}
