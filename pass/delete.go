package pass

import (
	"github.com/murtaza-u/z/pass/store"
	"github.com/rwxrob/bonzai/z"
)

var deleteCmd = &Z.Cmd{
	Name:    `delete`,
	Aliases: []string{"rm"},
	Summary: `delete an entry`,
	Usage:   `entry`,
	MinArgs: 1,
	Comp:    store.NewComp(),
	Call: func(caller *Z.Cmd, args ...string) error {
		d, err := Z.Conf.Data()
		if err != nil {
			return err
		}

		c, err := store.NewConfig([]byte(d))
		if err != nil {
			return err
		}
		s := store.New(c)

		return s.Delete(args...)
	},
}
