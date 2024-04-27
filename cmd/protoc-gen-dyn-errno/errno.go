package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func runProtoGen(gen *protogen.Plugin) error {
	if args.ErrorsPackage == "" {
		log.Fatal("errors package import path must be give with '--dyn-errno_out=epk=xxx'")
	}
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	for _, file := range gen.Files {
		if !file.Generate || len(file.Enums) == 0 {
			continue
		}
		enums := IntoEnums("", file.Enums)
		enums = append(enums, IntoEnumsFromMessage("", file.Messages)...)
		if len(enums) == 0 {
			continue
		}

		filename := file.GeneratedFilenamePrefix + ".errno.pb.go"
		g := gen.NewGeneratedFile(filename, file.GoImportPath)
		mt := &ErrnoFile{
			Version:       version,
			ProtocVersion: protocVersion(gen),
			IsDeprecated:  file.Proto.GetOptions().GetDeprecated(),
			Source:        file.Desc.Path(),
			Package:       string(file.GoPackageName),
			Epk:           args.ErrorsPackage,
			Errors:        enums,
		}
		err := mt.execute(g)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr,
				"\u001B[31mWARN\u001B[m: execute template failed. %v\n", err)
		}
	}
	return nil
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}
