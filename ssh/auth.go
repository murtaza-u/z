package ssh

import (
	"fmt"
	"os"

	"github.com/seehuhn/password"
	"golang.org/x/crypto/ssh"
)

func Auth(privF string) ([]ssh.AuthMethod, error) {
	if privF == "" {
		pswd, err := password.Read("passwd: ")
		if err != nil {
			return nil, fmt.Errorf("unable to get password: %w", err)
		}
		return []ssh.AuthMethod{
			ssh.Password(string(pswd)),
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
