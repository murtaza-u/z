package pass

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/structs/set/text/set"
)

func newComp() *comp {
	return new(comp)
}

type comp struct{}

func (comp) Complete(_ bonzai.Command, args ...string) []string {
	entries := listEntries()
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
