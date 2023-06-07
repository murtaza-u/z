package totp

import (
	"fmt"
	"time"

	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/z/pass/store"

	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

const warning = "otp copied to clipboard and will be cleared in 30s. Do *not* exit."

var copyCmd = &cli.Command{
	Name:         "copy",
	Aliases:      []string{"cp"},
	Usage:        "generate and copy otp to clipboard",
	UsageText:    "copy ENTRY",
	BashComplete: store.Comp,
	Action: func(ctx *cli.Context) error {
		conf := conf.New()
		conf.MustInit()

		d, err := conf.Data()
		if err != nil {
			return err
		}

		c, err := store.NewConfig([]byte(d))
		if err != nil {
			return err
		}
		s := store.New(c)

		out, err := s.Decrypt(ctx.Args().First())
		if err != nil {
			return err
		}

		otp, err := GenOTP(string(out))
		if err != nil {
			return fmt.Errorf("failed to generate TOTP: %w", err)
		}

		err = clipboard.Init()
		if err != nil {
			return fmt.Errorf("failed to initialise clipboard: %w", err)
		}

		changed := clipboard.Write(clipboard.FmtText, []byte(otp))
		fmt.Println(warning)

		t := time.NewTimer(time.Second * 30)

		select {
		case <-changed:
		case <-t.C:
			clipboard.Write(clipboard.FmtText, []byte{})
		}

		return nil
	},
}
