package vi

import (
	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
)

var Cmd = &Z.Cmd{
	Name:    `vi`,
	Summary: `neatvi embedded inside bonzai commander`,
	Comp:    compfile.New(),
	Call: func(caller *Z.Cmd, args ...string) error {
		return vi(args...)
	},
}
