package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bmaupin/go-epub"
	"github.com/pterm/pterm"

	"github.com/liberty239/cassiopaea-tools/pkg/source/adapter"
	"github.com/liberty239/cassiopaea-tools/pkg/source/local"
	"github.com/liberty239/cassiopaea-tools/pkg/term"
)

const epubCss = `
html {
	font-size: 14px
}
body {
	color: black;
	font-family: Georgia, Times, "Times New Roman", serif;
	font-size: 1.1rem;
	line-height: 1.5;
	text-align: justify;
}
h2 {
	color: black;
	font-family: Georgia, Times, "Times New Roman", serif;
	font-size: 2.2rem;
	font-weight: bold;
	line-height: 1.5;
	text-align: justify;
}
.answer {
	background-color: rgb(232,232,232);
}
`

func writeCss() (*os.File, error) {
	f, err := os.CreateTemp(os.TempDir(), "*.css")
	if err != nil {
		return nil, err
	}

	if _, err := f.WriteString(epubCss); err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}

func doEpub(args args) error {
	then := time.Now()

	s, err := term.Spinner().Start("Rendering...")
	if err != nil {
		s.Fail(err)
		return err
	}

	f, err := writeCss()
	if err != nil {
		s.Fail(err)
		return err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	e := epub.NewEpub("")
	css, err := e.AddCSS(f.Name(), "")
	if err != nil {
		return err
	}

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
			css,
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
