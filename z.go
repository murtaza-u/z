package z

import (
	"github.com/murtaza-u/z/age"
	"github.com/murtaza-u/z/isosec"
	"github.com/murtaza-u/z/pass"
	"github.com/murtaza-u/z/pomo"
	"github.com/murtaza-u/z/ssh"
	"github.com/murtaza-u/z/vi"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/yq"
)

var Cmd = &Z.Cmd{
	Name:      `z`,
	Summary:   `My personal stateful monolith Bonzai™ commander`,
	Copyright: `Copyright 2023 Murtaza Udaipurwala`,
	License:   `Apache-2.0`,
	Site:      `https://murtazau.xyz`,
	Source:    `https://github.com/murtaza-u/z`,
	Issues:    `https://github.com/murtaza-u/z/issues`,
	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, yq.Cmd, pomo.Cmd, vi.Cmd, isosec.Cmd,
		ssh.Cmd, age.Cmd, pass.Cmd,
	},
}

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
}
