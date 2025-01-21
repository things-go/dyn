package main

import (
	"os"

	"github.com/things-go/dyn/cmd/dyn-gen/command"
	_ "github.com/things-go/ens/driver/mysql"
)

func main() {
	err := command.NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
