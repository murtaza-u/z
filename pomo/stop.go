package pomo

import "github.com/rwxrob/bonzai/z"

var stopCmd = &Z.Cmd{
	Name:    `stop`,
	Summary: `stop the countdown timer`,
	NoArgs:  true,
	Call: func(caller *Z.Cmd, args ...string) error {
		Z.Vars.Del(".pomo.endt")
		return nil
	},
}
