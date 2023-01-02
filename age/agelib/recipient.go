package agelib

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/cryptobyte"
)

func GetRecipients(arg string) ([]age.Recipient, error) {
	var recs []age.Recipient

	if _, err := os.Stat(arg); err == nil {
		_recs, err := parseRecipientFile(arg)
		if err != nil {
			return nil, err
		}
		recs = append(recs, _recs...)

		return recs, nil
	}

	for _, arg := range strings.Split(arg, ",") {
		rec, err := parseRecipient(arg)
		if err != nil {
			return nil, err
		}

		recs = append(recs, rec)
	}

	return recs, nil
}

func parseRecipient(arg string) (age.Recipient, error) {
	if arg == "" {
		return nil, fmt.Errorf("missing recipient")
	}

	if strings.HasPrefix(arg, "age1") {
		return age.ParseX25519Recipient(arg)
	}

	if strings.HasPrefix(arg, "ssh-") {
		return agessh.ParseRecipient(arg)
	}

	return nil, fmt.Errorf("unknown recipient type: %q", arg)
}

func parseRecipientFile(name string) ([]age.Recipient, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open recipient file: %w", err)
	}
	defer f.Close()

	var recs []age.Recipient
	scanner := bufio.NewScanner(
		io.LimitReader(f, RecipientFileSizeLimit),
	)

	var n int

	for scanner.Scan() {
		n++

		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if len(line) > LineLengthLimit {
			return nil, fmt.Errorf("%q: line %d is too long", name, n)
		}

		r, err := parseRecipient(line)

		if err != nil {
			if t, ok := sshKeyType(line); ok {
				// skip unsupported but valid SSH public keys with a
				// warning.
				fmt.Fprintf(
					os.Stderr,
					"file %q: ignoring unsupported SSH key of type %q at line %d\n",
					name, t, n,
				)

				continue
			}

			// hide the error since it might unintentionally leak the
			// contents of confidential files.
			return nil, fmt.Errorf(
				"%q: malformed recipient at line %d", name, n,
			)
		}

		recs = append(recs, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%q: failed to read recipient file: %v", name, err)
	}

	if len(recs) == 0 {
		return nil, fmt.Errorf("%q: no recipients found", name)
	}

	return recs, nil
}

func sshKeyType(s string) (string, bool) {
	fields := strings.Split(s, " ")
	if len(fields) < 2 {
		return "", false
	}

	key, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		return "", false
	}

	var typLen uint32
	var typBytes []byte

	k := cryptobyte.String(key)
	if !k.ReadUint32(&typLen) || !k.ReadBytes(&typBytes, int(typLen)) {
		return "", false
	}

	if t := fields[0]; t == string(typBytes) {
		return t, true
	}

	return "", false
}
