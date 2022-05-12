package main

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	"github.com/things-go/dyn/errors"
)

const (
	fmtPackage = protogen.GoImportPath("fmt")
)

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
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.QualifiedGoIdent(fmtPackage.Ident(""))
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
			CamelValue: camelCase(string(v.Desc.Name())),
			Message:    msg,
		}
		ew.Errors = append(ew.Errors, err)
	}
	if len(ew.Errors) == 0 {
		return true
	}
	g.P(ew.execute())
	return false
}

func camelCase(str string) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}

	var b strings.Builder
	var words []string

	for i, s := 0, str; s != ""; s = s[i:] { // split on upper letter or _
		i = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if i <= 0 {
			i = len(s)
		}
		word := s[:i]
		words = append(words, strings.Split(word, "_")...)
	}

	for i, word := range words {
		word = removeInvalidChars(word, i == 0) // on 0 remove first digits
		if word == "" {
			continue
		}

		out := strings.ToUpper(string(word[0]))
		if len(word) > 1 {
			out += strings.ToLower(word[1:])
		}
		b.WriteString(out)
	}

	if b.Len() == 0 { // check if this is number
		if _, err := strconv.Atoi(str); err == nil {
			b.WriteString("Key")
			b.WriteString(str)
		}
	}

	return b.String()
}

func removeInvalidChars(s string, removeFirstDigit bool) string {
	var buf bytes.Buffer

	for _, b := range []byte(s) {
		if b >= 97 && b <= 122 { // a-z
			buf.WriteByte(b)
			continue
		}
		if b >= 65 && b <= 90 { // A-Z
			buf.WriteByte(b)
			continue
		}
		if b >= 48 && b <= 57 { // 0-9
			if !removeFirstDigit || buf.Len() > 0 {
				buf.WriteByte(b)
				continue
			}
		}
	}

	return buf.String()
}
