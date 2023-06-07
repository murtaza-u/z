package pass

import (
	"fmt"
	"strings"

	"github.com/murtaza-u/z/pass/store"

	"github.com/murtaza-u/conf"
	"github.com/urfave/cli/v2"
)

var showCmd = &cli.Command{
	Name:         "show",
	Aliases:      []string{"ls", "list"},
	Usage:        "list entries / decrypt an entry",
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
		fmt.Println(strings.TrimSpace(string(out)))

		return nil
	},
}
