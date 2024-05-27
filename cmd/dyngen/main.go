package main

import (
	"os"

	"github.com/things-go/dyn/cmd/dyngen/command"
)

var root = command.NewRootCmd()

func main() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
