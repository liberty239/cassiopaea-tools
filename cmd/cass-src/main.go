package main

import (
	"os"

	"github.com/docopt/docopt.go"

	"github.com/liberty239/cassiopaea-tools/pkg/source"
	"github.com/liberty239/cassiopaea-tools/pkg/source/forum"
	"github.com/liberty239/cassiopaea-tools/pkg/source/local"
)

var usage = `Cassiopaea Session Transcripts source.

Usage:
  cass-src list [--local=<directory>]
  cass-src fetch <date> <file> [--local=<directory>]
  cass-src fetch-all <directory> [--local=<directory>]
  cass-src -h | --help
  cass-src --version

Options:
  -h --help  Show this screen.
  --version  Show version.`

type args struct {
	List      bool   `docopt:"list"`
	Fetch     bool   `docopt:"fetch"`
	FetchAll  bool   `docopt:"fetch-all"`
	Date      string `docopt:"<date>"`
	File      string `docopt:"<file>"`
	Directory string `docopt:"<directory>"`
	Local     string `docopt:"--local"`
}

func argsToSource(args args) source.Source {
	if len(args.Local) > 0 {
		return local.NewSource(args.Local)
	}
	return forum.NewSource()
}

func main() {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}

	var args args
	if err := opts.Bind(&args); err != nil {
		panic(err)
	}

	switch {
	case args.List:
		err = doList(args)
	case args.Fetch:
		err = doFetch(args)
	case args.FetchAll:
		err = doFetchAll(args)
	default:
		panic("unknown command")
	}

	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
