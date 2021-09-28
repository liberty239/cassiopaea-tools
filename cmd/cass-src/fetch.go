package main

import (
	"io"
	"os"
	"time"

	"github.com/carmo-evan/strtotime"
	"github.com/pterm/pterm"

	"github.com/liberty239/cassiopaea-tools/pkg/term"
)

func doFetch(args args) error {
	then := time.Now()

	s, err := term.Spinner().Start("Fetching...")
	if err != nil {
		s.Fail(err)
		return err
	}

	ts, err := strtotime.Parse(args.Date, time.Now().Unix())
	if err != nil {
		s.Fail(err)
		return err
	}

	r, _, err := argsToSource(args).Fetch(time.Unix(ts, 0))
	if err != nil {
		s.Fail(err)
		return err
	}
	if r == nil {
		s.Fail("Session not found")
		return err
	}
	defer r.Close()

	f, err := os.Create(args.File)
	if err != nil {
		s.Fail(err)
		return err
	}

	if _, err := io.Copy(f, r); err != nil {
		s.Fail(err)
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		s.Fail(err)
		return err
	}

	s.Success(time.Unix(ts, 0).Format("02 January 2006"))
	pterm.Info.Println("Fetch took", time.Since(then).Seconds(), "seconds")

	return nil
}
