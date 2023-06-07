package isosec

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

const Format = "20060102150405"

var Cmd = &cli.Command{
	Name:        "isosec",
	Usage:       "create/parse isosec identifiers (YYYYMMDDhhmmss)",
	Subcommands: []*cli.Command{parseCmd},
	Action: func(ctx *cli.Context) error {
		t := time.Now().UTC().Format(Format)
		fmt.Println(t)
		return nil
	},
}

var parseCmd = &cli.Command{
	Name:      "parse",
	Usage:     "parse isosec indentifier",
	UsageText: "isosec-id",
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args().First()
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
