package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/things-go/dyn/genproto/errors"
	"github.com/things-go/dyn/internal/infra"
)

const (
	fmtPackage = protogen.GoImportPath("fmt")
)

func runProtoGen(gen *protogen.Plugin) error {
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
}

// generateFile generates a _error.pb.go file containing errors definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Enums) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_error.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-errno. DO NOT EDIT.")
	g.P("// versions:")
	g.P("//   - protoc-gen-go-errno ", version)
	g.P("//   - protoc              ", infra.ProtocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = ", fmtPackage.Ident("Errorf"))
	g.P()

	generateFileContent(gen, file, g)
	return g
}

// generateFileContent generates the errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Enums) == 0 {
		return
	}

	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the errors package it is being compiled against.")
	g.P("const _ = ", protogen.GoImportPath(*errorsPackage).Ident("ErrorsProtoPackageIsVersion3"))
	g.P()
	index := 0
	for _, enum := range file.Enums {
		skip := genErrorsDetail(gen, file, g, enum)
		if !skip {
			index++
		}
	}
	// If all enums do not contain 'errors.code', the current file is skipped
	if index == 0 {
		g.Skip()
	}
}

func genErrorsDetail(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, enum *protogen.Enum) bool {
	defaultCode := proto.GetExtension(enum.Desc.Options(), errors.E_DefaultCode)
	code := 0
	if ok := defaultCode.(int32); ok != 0 {
		code = int(ok)
	}
	var ew errorWrapper
	for _, v := range enum.Values {
		enumCode := code
		eCode := proto.GetExtension(v.Desc.Options(), errors.E_Code)
		if vv, ok := eCode.(int32); ok && vv != 0 {
			enumCode = int(vv)
		}
		if enumCode == 0 {
			continue
		}
		msg := ""
		eMsg := proto.GetExtension(v.Desc.Options(), errors.E_Msg)
		if vv, ok := eMsg.(string); ok {
			msg = vv
		}

		err := &errorInfo{
			Name:       string(enum.Desc.Name()),
			Code:       enumCode,
			Value:      string(v.Desc.Name()),
			CamelValue: infra.CamelCase(string(v.Desc.Name())),
			Message:    msg,
		}
		ew.Errors = append(ew.Errors, err)
	}
	if len(ew.Errors) == 0 {
		return true
	}
	err := ew.execute(g)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"\u001B[31mWARN\u001B[m: execute template failed.\n")
	}
	return false
}
