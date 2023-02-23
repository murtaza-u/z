package totp

import (
	"fmt"

	"github.com/murtaza-u/z/pass/store"

	Z "github.com/rwxrob/bonzai/z"
)

var showCmd = &Z.Cmd{
	Name:    `show`,
	Aliases: []string{"ls", "list"},
	Summary: `list entries / generate an otp`,
	Usage:   `[entry]`,
	MaxArgs: 1,
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

		if len(args) == 0 || args[0] == "" {
			for _, l := range s.List() {
				fmt.Println(l)
			}

			return nil
		}

		out, err := s.Decrypt(args[0])
		if err != nil {
			return err
		}

		otp, err := GenOTP(string(out))
		if err != nil {
			return fmt.Errorf("failed to generate TOTP: %w", err)
		}
		fmt.Println(otp)

		return nil
	},
}
