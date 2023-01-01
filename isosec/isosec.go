package isosec

import (
	"fmt"
	"time"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

const Format = "20060102150405"

var Cmd = &Z.Cmd{
	Name:     `isosec`,
	Summary:  `create/parse isosec identifiers (YYYYMMDDhhmmss)`,
	Commands: []*Z.Cmd{help.Cmd, parseCmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		t := time.Now().UTC().Format(Format)
		fmt.Println(t)
		return nil
	},
}

var parseCmd = &Z.Cmd{
	Name:    `parse`,
	Summary: `parse isosec indentifier`,
	NumArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		arg := args[0]
		if arg == "" {
			return fmt.Errorf("missing argument")
		}

		t, err := time.Parse(Format, arg)
		if err != nil {
			return fmt.Errorf("invalid isosec id %s", arg)
		}

		fmt.Printf("UTC: %v\nLocal: %v\n", t, t.Local())

		return nil
	},
}
