package totp

import (
	"fmt"
	"time"

	"github.com/murtaza-u/z/pass/store"

	"github.com/rwxrob/bonzai/z"
	"golang.design/x/clipboard"
)

const warning = "otp copied to clipboard and will be cleared in 30s. Do *not* exit."

var copyCmd = &Z.Cmd{
	Name:    `copy`,
	Aliases: []string{"cp"},
	Summary: `generate and copy otp to clipboard`,
	Usage:   `entry`,
	NumArgs: 1,
	Comp:    store.NewComp(),
	Call: func(caller *Z.Cmd, args ...string) error {
		d, err := Z.Conf.Data()
		if err != nil {
			return err
		}

		c, err := store.NewConfig([]byte(d))
		if err != nil {
			return err
		}
		s := store.New(c)

		out, err := s.Decrypt(args[0])
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
