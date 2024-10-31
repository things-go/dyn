package main

import (
	"flag"
	"fmt"

	"github.com/things-go/dyn/cmd/internal/meta"
	"google.golang.org/protobuf/compiler/protogen"
)

var args = struct {
	ShowVersion         bool
	Omitempty           bool
	AllowDeleteBody     bool
	AllowEmptyPatchBody bool
	UseEncoding         bool
	EnableMetadata      bool
}{
	ShowVersion:         false,
	Omitempty:           true,
	AllowDeleteBody:     false,
	AllowEmptyPatchBody: false,
	UseEncoding:         false,
	EnableMetadata:      false,
}

func init() {
	flag.BoolVar(&args.ShowVersion, "version", false, "print the version and exit")
	flag.BoolVar(&args.Omitempty, "omitempty", true, "omit if google.api is empty")
	flag.BoolVar(&args.AllowDeleteBody, "allow_delete_body", false, "allow delete body")
	flag.BoolVar(&args.AllowEmptyPatchBody, "allow_empty_patch_body", false, "allow empty patch body")
	flag.BoolVar(&args.UseEncoding, "use_encoding", false, "use the framework encoding")
	flag.BoolVar(&args.EnableMetadata, "enable_metadata", false, "store the metadata for every router.")
}

func main() {
	flag.Parse()
	if args.ShowVersion {
		fmt.Printf("protoc-gen-dyn-gin %v\n", meta.Version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
