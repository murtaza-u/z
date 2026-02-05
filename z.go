package z

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/murtaza-u/z/internal/vars"
	"github.com/murtaza-u/z/pomo"

	"github.com/urfave/cli/v3"
)

// Run initializes and executes the monolith commander based on the provided
// arguments.
func Run(ctx context.Context, args ...string) error {
	err := vars.New().Init()
	if err != nil {
		return fmt.Errorf("failed to initialize cache vars: %w", err)
	}
	cmd := cli.Command{
		Name:                  "z",
		Usage:                 "Go monolith commander",
		Version:               "v0.1.1",
		EnableShellCompletion: true,
		Copyright:             "Apache-2.0",
		Authors: []any{
			mail.Address{Name: "Murtaza Udaipurwala", Address: "murtaza@murtazau.xyz"},
		},
		Commands: []*cli.Command{pomo.Cmd},
	}
	return cmd.Run(ctx, args)
}
