package pomo

import (
	"github.com/murtaza-u/z/internal/vars"

	"github.com/urfave/cli/v2"
)

var stopCmd = &cli.Command{
	Name:  "stop",
	Usage: "stop the countdown timer",
	Action: func(ctx *cli.Context) error {
		vars := vars.New()
		vars.Init()
		vars.Del(".pomo.endt")
		return nil
	},
}
