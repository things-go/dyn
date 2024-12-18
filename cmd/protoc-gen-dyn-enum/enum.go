package main

import (
	"errors"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/things-go/dyn/cmd/internal/meta"
	"github.com/things-go/dyn/cmd/internal/protoenum"
	"github.com/things-go/dyn/cmd/internal/protoutil"
)

func runProtoGen(gen *protogen.Plugin) error {
	var mergeEnums []*protoenum.Enum
	var source []string
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	isMerge := args.Merge
	if isMerge {
		if args.Package == "" ||
			args.Filename == "" ||
			args.GoPackage == "" {
			return errors.New("when enable merge, filename,package,go_package must be set")
		}
		mergeEnums = make([]*protoenum.Enum, 0, len(gen.Files)*4)
		source = make([]string, 0, len(gen.Files))
	}
	usedTemplate := enumTemplate
	if args.CustomTemplate != "" {
		t, err := ParseTemplateFromFile(args.CustomTemplate)
		if err != nil {
			return err
		}
		usedTemplate = t
	}

	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		enums := protoenum.IntoEnums("", f.Enums)
		enums = append(enums, protoenum.IntoEnumsFromMessage("", f.Messages)...)
		if len(enums) == 0 {
			continue
		}
		if isMerge {
			source = append(source, f.Desc.Path())
			mergeEnums = append(mergeEnums, enums...)
			continue
		}
		g := gen.NewGeneratedFile(f.GeneratedFilenamePrefix+args.Suffix, f.GoImportPath)
		e := &EnumFile{
			Version:       meta.Version,
			ProtocVersion: protoutil.ProtocVersion(gen),
			IsDeprecated:  f.Proto.GetOptions().GetDeprecated(),
			Source:        f.Desc.Path(),
			Package:       string(f.GoPackageName),
			Enums:         enums,
		}
		_ = e.execute(usedTemplate, g)
	}
	if isMerge {
		g := gen.NewGeneratedFile(args.Filename+args.Suffix, protogen.GoImportPath(args.GoPackage))
		mergeFile := &EnumFile{
			Version:       meta.Version,
			ProtocVersion: protoutil.ProtocVersion(gen),
			IsDeprecated:  false,
			Source:        strings.Join(source, ","),
			Package:       args.Package,
			Enums:         mergeEnums,
		}
		return mergeFile.execute(usedTemplate, g)
	}
	return nil
}
