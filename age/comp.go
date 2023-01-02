package age

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/compcmd"
	"github.com/rwxrob/compfile"
)

func newComp() *comp {
	return new(comp)
}

type comp struct{}

func (comp) Complete(x bonzai.Command, args ...string) []string {
	if len(args) <= 1 {
		return compfile.New().Complete(x, args...)
	}

	return compcmd.New().Complete(x, args...)
}
