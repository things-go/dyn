package main

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "v0.0.2"

var showVersion = flag.Bool("version", false, "print the version and exit")
var errorsPackage = flag.String("epk", "github.com/things-go/dyn/errors", "errors core package in your project")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-errno %v\n", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if *errorsPackage == "" {
			log.Fatal("errors package import path must be give with '--go-errno_out=epk=xxx'")
		}
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}
