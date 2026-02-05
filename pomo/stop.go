package pomo

import (
	"context"

	"github.com/murtaza-u/z/internal/vars"

	"github.com/urfave/cli/v3"
)

var stopCmd = &cli.Command{
	Name:  "stop",
	Usage: "stop the countdown timer",
	Action: func(ctx context.Context, c *cli.Command) error {
		vars := vars.New()
		vars.Init()
		vars.Del(".pomo.endt")
		return nil
	},
}
