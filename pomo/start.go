package pomo

import (
	"fmt"
	"time"

	"github.com/murtaza-u/z/internal/vars"

	"github.com/urfave/cli/v2"
)

var startCmd = &cli.Command{
	Name:      "start",
	Usage:     "start the countdown timer",
	UsageText: "start [duration]",
	Action: func(ctx *cli.Context) error {
		vars := vars.New()
		_dur := ctx.Args().First()
		if _dur == "" {
			_dur = DefaultDuration
		}

		dur, err := time.ParseDuration(_dur)
		if err != nil {
			return fmt.Errorf("invalid duration %s: %w", _dur, err)
		}

		endt := time.Now().Add(dur).Format(time.RFC3339)
		return vars.Set(".pomo.endt", endt)
	},
}
