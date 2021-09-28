package main

import (
	"fmt"
	"io"
	"time"

	"github.com/bmaupin/go-epub"
	"github.com/pterm/pterm"

	"github.com/liberty239/cassiopaea-tools/pkg/source/adapter"
	"github.com/liberty239/cassiopaea-tools/pkg/source/local"
	"github.com/liberty239/cassiopaea-tools/pkg/term"
)

func doEpub(args args) error {
	then := time.Now()

	s, err := term.Spinner().Start("Rendering...")
	if err != nil {
		s.Fail(err)
		return err
	}

	e := epub.NewEpub("")

	from, to := time.Now(), time.Time{}

	src := adapter.Chain(
		local.NewSource(args.Directory),
		adapter.Ascending,
		adapter.Synchronized,
	)
	if err := src.FetchAll(func(ts time.Time, url string, rc io.ReadCloser) error {
		b, err := io.ReadAll(rc)
		if err != nil {
			rc.Close()
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}

		if _, err := e.AddSection(
			string(b),
			ts.Format("02 January 2006"),
			ts.Format("2006-01-02.xhtml"),
			"",
		); err != nil {
			return err
		}

		switch {
		case ts.Before(from):
			from = ts
		case ts.After(to):
			to = ts
		}

		return nil
	}); err != nil {
		s.Fail(err)
		return err
	}

	e.SetAuthor("cassiopaea.org")
	e.SetTitle(
		fmt.Sprintf(
			"Cassiopaea Session Transcripts %s - %s",
			from.Format("02 January 2006"),
			to.Format("02 January 2006"),
		),
	)

	if err := e.Write(args.File); err != nil {
		s.Fail(err)
		return err
	}

	s.Success("Sessions written to EPUB file: ", args.File)
	pterm.Info.Println("Rendering took", time.Since(then).Seconds(), "seconds")

	return nil
}
