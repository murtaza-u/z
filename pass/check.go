package pass

import (
	"fmt"

	"github.com/rwxrob/bonzai/z"
)

var checkCmd = &Z.Cmd{
	Name:    `check`,
	Summary: `check if pass is setup correctly`,
	NoArgs:  true,
	Call: func(caller *Z.Cmd, args ...string) error {
		_, err := newCfg()
		if err != nil {
			return err
		}

		fmt.Println("everything looks great!")

		return nil
	},
}
