package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "v0.1.2"

var showVersion = flag.Bool("version", false, "print the version and exit")
var omitempty = flag.Bool("omitempty", true, "omit if google.api is empty")
var allowDeleteBody = flag.Bool("allow_delete_body", false, "allow delete body")
var allowEmptyPatchBody = flag.Bool("allow_empty_patch_body", false, "allow empty patch body")
var useCustomResponse = flag.Bool("use_custom_response", false, "use custom response encoder")
var rpcMode = flag.String("rpc_mode", "", "rpc mode, default empty use official rpc, options: rpcx")
var allowFromAPI = flag.Bool("allow_from_api", false, "allow from api can convert different api format.")
var useEncoding = flag.Bool("use_encoding", false, "use the framework encoding")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-gin %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
