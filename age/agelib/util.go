package agelib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/seehuhn/password"
)

func ReadIn(fname string) (io.Reader, error) {
	if fname == "" {
		fname = os.Stdin.Name()
	}

	f, err := os.OpenFile(fname, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open file %q: %w", fname, err,
		)
	}

	return bufio.NewReader(f), nil
}

func OpenOut(fname string) (io.WriteCloser, error) {
	if fname == "" {
		return os.Stdout, nil
	}

	out, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open file %q: %w", fname, err,
		)
	}

	return out, nil
}

func ReadHidden(form string, args ...string) string {
	prompt := strings.Join(append([]string{form}, args...), " ")

	var pswd []byte
	for pswd == nil || len(pswd) == 0 {
		pswd, _ = password.Read(prompt)
	}

	return string(pswd)
}
