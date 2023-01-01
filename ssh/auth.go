package ssh

import (
	"fmt"
	"os"

	"github.com/rwxrob/term"
	"golang.org/x/crypto/ssh"
)

func Auth(privF string) ([]ssh.AuthMethod, error) {
	if privF == "" {
		var pswd string

		fmt.Print("Enter password: ")
		for pswd == "" {
			pswd = term.ReadHidden()
		}

		return []ssh.AuthMethod{
			ssh.Password(pswd),
		}, nil
	}

	key, err := os.ReadFile(privF)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file %q: %w", privF, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}, nil
}
