package main

import (
	"go.uber.org/zap"

	"github.com/things-go/dyn/zapl"
)

func main() {
	l := zapl.New(zapl.WithLevel("debug"))
	zapl.ReplaceGlobals(l.Sugar().With(zap.String("main", "debug")))

	zapl.Debug("hello")
}
