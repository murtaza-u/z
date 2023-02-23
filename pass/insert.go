package pass

import (
	"github.com/murtaza-u/z/pass/store"

	"github.com/rwxrob/bonzai/z"
)

var insertCmd = &Z.Cmd{
	Name:    `insert`,
	Summary: `insert a new password entry`,
	NumArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		d, err := Z.Conf.Data()
		if err != nil {
			return err
		}

		c, err := store.NewConfig([]byte(d))
		if err != nil {
			return err
		}

		e := args[0]

		s := store.New(c)
		out, err := s.Insert(e)
		if err != nil {
			return err
		}

		return s.WriteEntry(e, out)
	},
}
