package main

import (
	"log"

	"github.com/murtaza-u/z/ssh"
)

func main() {
	log.SetFlags(0)
	ssh.Cmd.Run()
}
