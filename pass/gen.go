package pass

import (
	"fmt"
	"strconv"

	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli/v2"
)

const defaultLength = 50

var genCmd = &cli.Command{
	Name:      "gen",
	Usage:     "generate password",
	UsageText: "gen [LENGTH]",
	Action: func(ctx *cli.Context) error {
		n := defaultLength

		if l := ctx.Args().First(); l != "" {
			conv, err := strconv.Atoi(l)
			if err != nil {
				return fmt.Errorf("failed to parse argument: %w", err)
			}
			if conv <= 0 {
				return fmt.Errorf("%d: length must be > 0", conv)
			}
			n = conv
		}

		pswd, err := password.Generate(n, 10, 10, false, true)
		if err != nil {
			return err
		}
		fmt.Println(pswd)

		return nil
	},
}
