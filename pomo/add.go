package pomo

import (
	"fmt"
	"time"

	"github.com/murtaza-u/z/internal/vars"

	"github.com/urfave/cli/v2"
)

var addCmd = &cli.Command{
	Name:        "add",
	Usage:       "extend duration of on-going countdown timer",
	UsageText:   "add DURATION",
	Description: "Eg: add 10m5s",
	Action: func(ctx *cli.Context) error {
		vars := vars.New()
		vars.Init()

		_x := ctx.Args().First()
		if _x == "" {
			return fmt.Errorf("missing duration")
		}

		x, err := time.ParseDuration(_x)
		if err != nil {
			return fmt.Errorf("invalid duration %s: %w", x, err)
		}

		endt := vars.Get(".pomo.endt")
		if endt == "null" {
			return fmt.Errorf("countdown not started")
		}

		dur, err := time.Parse(time.RFC3339, endt)
		if err != nil {
			return fmt.Errorf(
				"error parsing duration %s: %w", endt, err,
			)
		}

		sub := time.Until(dur).Truncate(time.Second)
		if sub <= 0 {
			dur = time.Now()
		}

		dur = dur.Add(x)
		return vars.Set(".pomo.endt", dur.Format(time.RFC3339))
	},
}
