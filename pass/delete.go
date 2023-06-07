package pass

import (
	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/z/pass/store"
	"github.com/urfave/cli/v2"
)

var deleteCmd = &cli.Command{
	Name:         "delete",
	Aliases:      []string{"rm"},
	Usage:        "delete an entry",
	UsageText:    "delete ENTRY",
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

		return s.Delete(ctx.Args().Slice()...)
	},
}
