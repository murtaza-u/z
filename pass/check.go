package pass

import (
	"fmt"

	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/z/pass/store"

	"github.com/urfave/cli/v2"
)

var checkCmd = &cli.Command{
	Name:  "check",
	Usage: "check if pass is setup correctly",
	Action: func(ctx *cli.Context) error {
		conf := conf.New()
		conf.MustInit()

		d, err := conf.Data()
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
