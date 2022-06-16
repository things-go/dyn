package main

import (
	"github.com/things-go/dyn/log"
)

func main() {
	l := log.NewLoggerWith(log.New(log.WithLevel("debug")))
	log.ReplaceGlobals(l)

	log.Debug("hello")
}
