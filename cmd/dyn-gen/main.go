package main

import (
	"os"

	"github.com/things-go/dyn/cmd/dyn-gen/command"
)

func main() {
	err := command.NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
