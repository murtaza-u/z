package totp

import (
	"fmt"

	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/z/pass/store"

	"github.com/urfave/cli/v2"
)

var showCmd = &cli.Command{
	Name:         "show",
	Aliases:      []string{"ls", "list"},
	Usage:        "list entries / generate an otp",
	UsageText:    "show [ENTRY]",
	BashComplete: store.Comp,
	Action: func(ctx *cli.Context) error {
		conf := conf.New()
		conf.MustInit()

		d, err := conf.Data()
		if err != nil {
			return err
		}

		c, err := store.NewConfig([]byte(d))
		if err != nil {
			return err
		}
		s := store.New(c)

		arg := ctx.Args().First()
		if arg == "" {
			for _, l := range s.List() {
				fmt.Println(l)
			}
			return nil
		}

		out, err := s.Decrypt(arg)
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
