package pomo

import (
	"fmt"
	"time"

	"github.com/murtaza-u/z/internal/vars"

	"github.com/urfave/cli/v2"
)

const (
	DefaultDuration = "25m"
	Prefix          = "🍅"
	PrefixWarn      = "💢"
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

		endt := vars.Get(".pomo.endt")
		if endt == "" {
			return nil
		}

		dur, err := time.Parse(time.RFC3339, endt)
		if err != nil {
			return fmt.Errorf("error parsing duration %s: %w", endt, err)
		}

		sub := time.Until(dur).Truncate(time.Second)
		if sub > 0 {
			fmt.Printf("%s %s\n", Prefix, sub)
			return nil
		}
		fmt.Printf("%s %s\n", PrefixWarn, sub)

		return nil
	},
}
