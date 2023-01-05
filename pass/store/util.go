package store

import (
	"fmt"
	"strings"
)

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
