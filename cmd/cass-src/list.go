package main

import (
	"time"

	"github.com/pterm/pterm"

	"github.com/liberty239/cassiopaea-tools/pkg/term"
)

func doList(args args) error {
	then := time.Now()

	s, err := term.Spinner().Start("Listing...")
	if err != nil {
		s.Fail(err)
		return err
	}

	xs, err := argsToSource(args).List()
	if err != nil {
		s.Fail(err)
		return err
	}

	s.Success("Sessions available: ", len(xs))

	data := [][]string{
		{"Year", "Day and month"},
	}
	for _, x := range xs {
		data = append(data, []string{x.Format("2006"), x.Format("2 January")})
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithData(data).
		Render()

	pterm.Info.Println("Listing took", time.Since(then).Seconds(), "seconds")

	return nil
}
