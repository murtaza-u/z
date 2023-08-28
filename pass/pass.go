package pass

import (
	"github.com/murtaza-u/z/pass/totp"

	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "pass",
	Usage: "password manager based on AGE",
	Subcommands: []*cli.Command{
		showCmd, genCmd, checkCmd, insertCmd, deleteCmd, reEncryptCmd, totp.Cmd,
	},
}
