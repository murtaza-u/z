package pomo

import (
	"fmt"
	"time"

	"github.com/rwxrob/bonzai/z"
)

var addCmd = &Z.Cmd{
	Name:    `add`,
	Summary: `extend duration of on-going countdown timer`,
	NumArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		_x := args[0]
		if _x == "" {
			return fmt.Errorf("missing duration")
		}

		x, err := time.ParseDuration(_x)
		if err != nil {
			return fmt.Errorf("invalid duration %s: %w", x, err)
		}

		endt := Z.Vars.Get(".pomo.endt")
		if endt == "null" || endt == "" {
			return fmt.Errorf("countdown not started")
		}

		dur, err := time.Parse(time.RFC3339, endt)
		if err != nil {
			return fmt.Errorf(
				"error parsing duration %s: %w", endt, err,
			)
		}

		sub := dur.Sub(time.Now()).Truncate(time.Second)
		if sub <= 0 {
			dur = time.Now()
		}

		dur = dur.Add(x)

		return Z.Vars.Set(".pomo.endt", dur.Format(time.RFC3339))
	},
}
