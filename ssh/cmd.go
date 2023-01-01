package ssh

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
	"github.com/rwxrob/help"
)

var Cmd = &Z.Cmd{
	Name:     `ssh`,
	Summary:  `ssh client`,
	Usage:    `destination [command]`,
	Commands: []*Z.Cmd{help.Cmd},
	Keys: Z.Keys{
		{
			Name:  `key`,
			Usage: `path to private key`,
			Comp:  compfile.New(),
		},
	},
	MinArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		dst := args[0]
		if dst == "" {
			return fmt.Errorf("missing destination")
		}

		privF := caller.GetVal("key")
		method, err := Auth(privF)
		if err != nil {
			return err
		}

		c, err := NewClient(dst, method)
		if err != nil {
			return err
		}

		s, err := c.NewSession()
		if err != nil {
			return err
		}

		err = s.SetPipes()
		if err != nil {
			return err
		}

		go s.HandleSignals()

		if len(args) == 1 {
			err := s.Shell()
			if err != nil {
				fmt.Println(err.Error())
			}

			s.WaitAndClose()

			return nil
		}

		cmd := strings.Join(args[1:], " ")
		err = s.Run(cmd)
		if err != nil {
			fmt.Println(err.Error())
		}

		s.CloseAndWaith()

		return nil
	},
}
