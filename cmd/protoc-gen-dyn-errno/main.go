package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "v1.0.0"

var args = &struct {
	ShowVersion   bool   // 显示版本
	ErrorsPackage string // error package
}{

	ShowVersion:   false,
	ErrorsPackage: "",
}

func init() {
	flag.BoolVar(&args.ShowVersion, "version", false, "print the version and exit")
	flag.StringVar(&args.ErrorsPackage, "epk", "github.com/things-go/dyn/errorx", "errors core package in your project")
}

func main() {
	flag.Parse()
	if args.ShowVersion {
		fmt.Printf("protoc-gen-dyn-errno %v\n", version)
		return
	}
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
