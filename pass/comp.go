package pass

import (
	"github.com/murtaza-u/z/pass/store"
	"github.com/rwxrob/bonzai"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/structs/set/text/set"
)

func newComp() *comp {
	return new(comp)
}

type comp struct{}

func (comp) Complete(_ bonzai.Command, args ...string) []string {
	d, err := Z.Conf.Data()
	if err != nil {
		return nil
	}

	c, err := store.NewConfig([]byte(d), "")
	if err != nil {
		return nil
	}

	s := store.New(c)

	entries := s.List()
	if len(args) == 0 {
		return entries
	}

	nargs := len(args)

	// remove duplicates
	if nargs > 1 {
		entries = set.Minus[string, string](entries, args[:nargs-1])
	}

	return filt.HasPrefix(entries, args[nargs-1])
}
