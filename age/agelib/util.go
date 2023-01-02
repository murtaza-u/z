package agelib

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/rwxrob/term"
	"golang.org/x/crypto/ssh"
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

func readPubFile(name string) (ssh.PublicKey, error) {
	if name == "-" {
		return nil, fmt.Errorf(
			`failed to obtain public key for "-" SSH key Use a file for
			which the corresponding ".pub" file exists, or convert the
			private key to a modern format with "ssh-keygen -p -m
			RFC4716"`,
		)
	}

	f, err := os.Open(name + ".pub")
	if err != nil {
		return nil, fmt.Errorf(
			`failed to obtain public key for %q SSH key: %v Ensure %q
			exists, or convert the private key %q to a modern format
			with "ssh-keygen -p -m RFC4716"`,
			name, err, name+".pub", name,
		)
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read %q: %v", name+".pub", err,
		)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(contents)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse %q: %v", name+".pub", err,
		)
	}

	return pubKey, nil
}
