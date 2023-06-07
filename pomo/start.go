package pomo

import (
	"fmt"
	"time"

	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/conf/vars"
	"github.com/urfave/cli/v2"
)

var startCmd = &cli.Command{
	Name:      "start",
	Usage:     "start the countdown timer",
	UsageText: "start [duration]",
	Action: func(ctx *cli.Context) error {
		conf := conf.New()
		vars := vars.New()

		_dur, err := conf.Query(".pomo.duration")
		if err != nil || _dur == "null" {
			_dur = DefaultDuration
		}

		if arg := ctx.Args().First(); arg != "" {
			_dur = arg
		}

		dur, err := time.ParseDuration(_dur)
		if err != nil {
			return fmt.Errorf("invalid duration %s: %w", _dur, err)
		}

		endt := time.Now().Add(dur).Format(time.RFC3339)
		return vars.Set(".pomo.endt", endt)
	},
}
