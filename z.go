package z

import (
	"fmt"

	"github.com/murtaza-u/z/internal/vars"
	"github.com/murtaza-u/z/pomo"

	"github.com/urfave/cli/v2"
)

// Run initializes and executes the monolith commander based on the provided
// arguments.
func Run(args ...string) error {
	err := vars.New().Init()
	if err != nil {
		return fmt.Errorf("failed to initialize cache vars: %w", err)
	}

	app := cli.NewApp()
	app.Name = "z"
	app.Usage = "Go monolith commander"
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.Copyright = "Apache-2.0"
	app.Authors = []*cli.Author{
		{Name: "Murtaza Udaipurwala", Email: "murtaza@murtazau.xyz"},
	}
	app.Commands = []*cli.Command{pomo.Cmd}
	return app.Run(args)
}
