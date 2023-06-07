package store

import (
	"fmt"

	"github.com/murtaza-u/conf"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/structs/set/text/set"
	"github.com/urfave/cli/v2"
)

func Comp(ctx *cli.Context) {
	conf := conf.New()
	err := conf.Init()
	if err != nil {
		return
	}

	d, err := conf.Data()
	if err != nil {
		return
	}

	c, err := NewConfig([]byte(d))
	if err != nil {
		return
	}

	s := New(c)
	entries := s.List()
	nargs := ctx.NArg()
	if nargs == 0 {
		for _, e := range entries {
			fmt.Println(e)
		}
	}

	// remove duplicates
	if nargs > 1 {
		entries = set.Minus[string, string](entries, ctx.Args().Slice()[:nargs-1])
	}

	for _, e := range filt.HasPrefix(entries, ctx.Args().Get(nargs-1)) {
		fmt.Println(e)
	}
}
