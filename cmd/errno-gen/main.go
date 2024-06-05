package main

import (
	"os"

	"github.com/things-go/dyn/cmd/errno-gen/command"
)

var cmd = command.NewRootCmd()

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
