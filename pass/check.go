package pass

import (
	"fmt"

	"github.com/murtaza-u/z/pass/store"

	"github.com/rwxrob/bonzai/z"
)

var checkCmd = &Z.Cmd{
	Name:    `check`,
	Summary: `check if pass is setup correctly`,
	NoArgs:  true,
	Call: func(caller *Z.Cmd, args ...string) error {
		d, err := Z.Conf.Data()
		if err != nil {
			return err
		}

		_, err = store.NewConfig([]byte(d))
		if err != nil {
			return err
		}

		fmt.Println("everything looks great!")

		return nil
	},
}
