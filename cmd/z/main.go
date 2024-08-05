package main

import (
	"log"
	"os"

	"github.com/murtaza-u/z"
)

func main() {
	log.SetFlags(0)
	err := z.Run(os.Args...)
	if err != nil {
		log.Fatal(err)
	}
}
