package pass

import (
	"github.com/murtaza-u/z/pass/totp"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
)

var Cmd = &Z.Cmd{
	Name:    `pass`,
	Summary: `password manager based on AGE`,
	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, showCmd, checkCmd, insertCmd, deleteCmd,
		copyCmd, reencryptCmd, totp.Cmd,
	},
}

func init() {
	Z.Conf.SoftInit()
}
