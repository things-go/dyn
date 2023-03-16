package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "v0.1.2"

var showVersion = flag.Bool("version", false, "print the version and exit")
var errorsPackage = flag.String("epk", "github.com/things-go/dyn/genproto/errors", "errors core package in your project")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-errno %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
