package pomo

import (
	"fmt"
	"time"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

const (
	DefaultDuration   = "25m"
	DefaultPrefix     = "🍅"
	DefaultPrefixWarn = "💢"
)

var Cmd = &Z.Cmd{
	Name:    `pomo`,
	Summary: `sets or prints a countdown timer (with tomato)`,
	Description: `The Pomodoro technique is a simple yet effective tool
		for focused work with planned breaks in between. Francesco
		Cirillo coined the term "pomodoro" which translates to tomato,
		in the late 1980s after the tomato-shaped timer he used as a
		university student.`,
	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, vars.Cmd, startCmd, stopCmd, addCmd,
	},
	Call: func(caller *Z.Cmd, args ...string) error {
		endt := Z.Vars.Get(".pomo.endt")
		if endt == "null" {
			return nil
		}

		prefix, err := Z.Conf.Query(".pomo.prefix")
		if err != nil || prefix == "null" {
			prefix = DefaultPrefix
		}

		warn, err := Z.Conf.Query(".pomo.prefixWarn")
		if err != nil || warn == "null" {
			warn = DefaultPrefixWarn
		}

		dur, err := time.Parse(time.RFC3339, endt)
		if err != nil {
			return fmt.Errorf(
				"error parsing duration %s: %w", endt, err,
			)
		}

		sub := dur.Sub(time.Now()).Truncate(time.Second)
		if sub > 0 {
			fmt.Printf("%s %s\n", prefix, sub)
			return nil
		}
		fmt.Printf("%s %s\n", warn, sub)

		return nil
	},
}

func init() {
	Z.Vars.SoftInit()
	Z.Conf.SoftInit()
}
