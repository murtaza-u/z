package pomo

import (
	"fmt"
	"time"

	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/conf/vars"
	"github.com/urfave/cli/v2"
)

const (
	DefaultDuration   = "25m"
	DefaultPrefix     = "🍅"
	DefaultPrefixWarn = "💢"
)

var Cmd = &cli.Command{
	Name:      "pomo",
	Usage:     "sets or prints a countdown timer (with tomato)",
	UsageText: "pomo [start|stop|add] [arguments]",
	Description: `The Pomodoro technique is a simple yet effective tool
for focused work with planned breaks in between. Francesco Cirillo
coined the term "pomodoro" which translates to tomato, in the late
1980s after the tomato-shaped timer he used as a university
student.`,
	Subcommands: []*cli.Command{startCmd, stopCmd, addCmd},
	Action: func(ctx *cli.Context) error {
		vars := vars.New()
		vars.Init()
		conf := conf.New()

		endt := vars.Get(".pomo.endt")
		if endt == "" {
			return nil
		}

		prefix, err := conf.Query(".pomo.prefix")
		if err != nil || prefix == "null" {
			prefix = DefaultPrefix
		}

		warn, err := conf.Query(".pomo.prefixWarn")
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
