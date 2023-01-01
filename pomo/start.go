package pomo

import (
	"fmt"
	"time"

	"github.com/rwxrob/bonzai/z"
)

var startCmd = &Z.Cmd{
	Name:    `start`,
	Usage:   `duration(optional)`,
	Summary: `start the countdown timer`,
	MaxArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		_dur, err := Z.Conf.Query(".pomo.duration")
		if err != nil || _dur == "null" {
			_dur = DefaultDuration
		}

		if len(args) > 0 {
			_dur = args[0]
		}

		dur, err := time.ParseDuration(_dur)
		if err != nil {
			return fmt.Errorf("invalid duration %s: %w", _dur, err)
		}

		endt := time.Now().Add(dur).Format(time.RFC3339)

		return Z.Vars.Set(".pomo.endt", endt)
	},
}
