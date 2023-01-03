package pass

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func entryExists(entry string) (bool, error) {
	path := filepath.Join(Store, entry)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to stat %q: %w", path, err)
	}

	return false, nil
}

func listEntries() []string {
	var list []string

	filepath.WalkDir(Store, func(_ string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if e.IsDir() {
			return nil
		}

		list = append(list, e.Name())

		return nil
	})

	return list
}

var okays = []string{"y", "yes"}

func confirm(prompt string) bool {
	var inp string

	fmt.Printf("%s (y/N): ", prompt)
	fmt.Scanf("%s", &inp)
	inp = strings.TrimSpace(inp)

	for _, ok := range okays {
		if strings.EqualFold(inp, ok) {
			return true
		}
	}

	return false
}
