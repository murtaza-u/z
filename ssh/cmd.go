package ssh

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:      "ssh",
	Usage:     "ssh client",
	UsageText: "destination [--key FILE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "key",
			Usage:     "path to private key `FILE`",
			TakesFile: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		dst := ctx.Args().First()
		if dst == "" {
			return fmt.Errorf("missing destination")
		}

		privF := ctx.String("key")
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

		if ctx.NArg() == 1 {
			err := s.Shell()
			if err != nil {
				fmt.Println(err.Error())
			}
			s.WaitAndClose()
			return nil
		}

		cmd := strings.Join(ctx.Args().Tail(), " ")
		err = s.Run(cmd)
		if err != nil {
			fmt.Println(err.Error())
		}
		s.CloseAndWaith()
		return nil
	},
}
