//go:build !linux

package vi

import "fmt"

func vi(args ...string) error {
	return fmt.Errorf("OS not supported")
}
