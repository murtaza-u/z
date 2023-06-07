package z

import (
	"log"

	"github.com/murtaza-u/z/age"
	"github.com/murtaza-u/z/isosec"
	"github.com/murtaza-u/z/pass"
	"github.com/murtaza-u/z/pomo"
	"github.com/murtaza-u/z/ssh"

	"github.com/murtaza-u/conf"
	"github.com/murtaza-u/conf/vars"
	"github.com/urfave/cli/v2"
)

func Run(args ...string) error {
	app := cli.NewApp()
	app.Name = "z"
	app.Usage = "Go monolith commander"
	app.Version = "0.2"
	app.EnableBashCompletion = true
	app.Copyright = "Apache-2.0"
	app.Authors = []*cli.Author{
		{Name: "Murtaza Udaipurwala", Email: "murtaza@murtazau.xyz"},
	}
	app.Commands = []*cli.Command{
		pomo.Cmd, isosec.Cmd, ssh.Cmd, age.Cmd, pass.Cmd,
	}
	return app.Run(args)
}

func init() {
	err := vars.New().Init()
	if err != nil {
		log.Fatal(err)
	}

	err = conf.New().Init()
	if err != nil {
		log.Fatal(err)
	}
}
