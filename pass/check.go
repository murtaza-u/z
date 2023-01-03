package pass

import (
	"fmt"

	"github.com/murtaza-u/z/age/agelib"
	"github.com/rwxrob/bonzai/z"
)

var checkCmd = &Z.Cmd{
	Name:    `check`,
	Summary: `check if pass is setup correctly`,
	NoArgs:  true,
	Call: func(caller *Z.Cmd, args ...string) error {
		pub, err := Z.Conf.Query(".pass.pub")
		if err != nil {
			return err
		}

		if pub == "null" {
			return fmt.Errorf(
				".pass.pub (public key) not set in config",
			)
		}

		key, err := Z.Conf.Query(".pass.key")
		if err != nil {
			return err
		}

		if key == "null" {
			return fmt.Errorf(
				".pass.key (private key) not set in config",
			)
		}

		_, err = agelib.ParseIdentityFile(key)
		if err != nil {
			return fmt.Errorf("failed to parse key %q: %w", key, err)
		}

		fmt.Println("everything looks great!")

		return nil
	},
}
