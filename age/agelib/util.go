package agelib

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/rwxrob/term"
)

func ReadIn(fname string) (io.Reader, error) {
	if fname == "" {
		return nil, fmt.Errorf("missing file")
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

func ReadHidden(form string, args ...any) string {
	var pswd string

	for pswd == "" {
		pswd = term.PromptHidden(form, args...)
		fmt.Println()
	}

	return pswd
}
