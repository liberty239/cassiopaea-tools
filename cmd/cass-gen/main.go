package main

import (
	"os"

	"github.com/docopt/docopt.go"
)

var usage = `Cassiopaea Session Generator.

Usage:
  cass-gen epub <directory> <file>
  cass-gen html <directory> <file> [--meta-tags=<tags>]
  cass-gen -h | --help
  cass-gen --version

Options:
  --meta-tags=<tags>  Comma separated list of key:value pairs.
  -h --help           Show this screen.
  --version           Show version.`

type args struct {
	Epub      bool   `docopt:"epub"`
	Html      bool   `docopt:"html"`
	File      string `docopt:"<file>"`
	Directory string `docopt:"<directory>"`
	MetaTags  string `docopt:"--meta-tags"`
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
	case args.Epub:
		err = doEpub(args)
	case args.Html:
		err = doHtml(args)
	default:
		panic("unknown command")
	}

	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
