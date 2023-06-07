package pass

import (
	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/z/pass/store"

	"github.com/urfave/cli/v2"
)

var insertCmd = &cli.Command{
	Name:         "insert",
	Usage:        "insert a new password entry",
	UsageText:    "insert ENTRY",
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

		e := ctx.Args().First()
		s := store.New(c)
		out, err := s.Insert(e)
		if err != nil {
			return err
		}

		return s.WriteEntry(e, out)
	},
}
