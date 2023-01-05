package totp

import (
	"github.com/murtaza-u/z/age/agelib"
	"github.com/murtaza-u/z/pass/store"

	"github.com/rwxrob/bonzai/z"
)

var insertCmd = &Z.Cmd{
	Name:    `insert`,
	Summary: `insert new TOTP URI entry`,
	NumArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		d, err := Z.Conf.Data()
		if err != nil {
			return err
		}

		c, err := store.NewConfig([]byte(d), SubPath)
		if err != nil {
			return err
		}

		e := args[0]

		s := store.New(c)
		s.InputSecret = func() (string, error) {
			uri := agelib.ReadHidden("totp uri: ")
			return uri, nil
		}

		out, err := s.Insert(e)
		if err != nil {
			return err
		}

		return s.WriteEntry(e, out)
	},
}
