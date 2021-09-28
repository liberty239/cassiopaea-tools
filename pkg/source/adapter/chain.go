package adapter

import "github.com/liberty239/cassiopaea-tools/pkg/source"

func Chain(src source.Source, fs ...func(source.Source) source.Source) source.Source {
	for _, f := range fs {
		src = f(src)
	}
	return src
}
