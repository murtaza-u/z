package agelib

import (
	"bytes"
	"fmt"
	"os"

	"filippo.io/age"
	"github.com/rwxrob/fs/file"
)

func ParseRecipients(args ...string) ([]age.Recipient, error) {
	buf := new(bytes.Buffer)

	for _, arg := range args {
		if !file.Exists(arg) {
			_, err := buf.WriteString(arg + "\n")
			if err != nil {
				return nil, fmt.Errorf(
					"failed to write to buffer: %w", err,
				)
			}

			continue
		}

		data, err := os.ReadFile(arg)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to read file %q: %w", arg, err,
			)
		}

		data = append(data, []byte("\n")...)
		_, err = buf.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}
	}

	return age.ParseRecipients(buf)
}
