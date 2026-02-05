package pomo

import (
	"context"
	"fmt"
	"time"

	"github.com/murtaza-u/z/internal/vars"

	"github.com/urfave/cli/v3"
)

var startCmd = &cli.Command{
	Name:      "start",
	Usage:     "start the countdown timer",
	UsageText: "start [duration]",
	Action: func(ctx context.Context, c *cli.Command) error {
		vars := vars.New()
		_dur := c.Args().First()
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
