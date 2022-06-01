package main

import (
	"github.com/things-go/dyn/zapl"
)

func main() {
	l := zapl.NewLoggerWith(zapl.New(zapl.WithLevel("debug")))
	zapl.ReplaceGlobals(l)

	zapl.Debug("hello")
}
