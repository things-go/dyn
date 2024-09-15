package api

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/things-go/ens/proto"
	"github.com/things-go/ens/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const googleProtobufTimestamp = "google.protobuf.Timestamp"

type CodeGen struct {
	buf                       bytes.Buffer
	Entity                    *proto.Message    // required, proto Message
	PackageName               string            // required, proto 包名
	Options                   map[string]string // required, proto option
	Style                     string            // 字段代码风格, snakeCase, smallCamelCase, pascalCase
	DisableBool               bool              // 禁用bool,使用int32
	DisableTimestamp          bool              // 禁用google.protobuf.Timestamp,使用int64
	EnableOpenapiv2Annotation bool              // 启用int64的openapiv2注解
}

// Bytes returns the CodeBuf's buffer.
func (g *CodeGen) Bytes() []byte {
	return g.buf.Bytes()
}

// Write appends the contents of p to the buffer,
func (g *CodeGen) Write(b []byte) (n int, err error) {
	return g.buf.Write(b)
}

// Print formats using the default formats for its operands and writes to the generated output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Print(a ...any) (n int, err error) {
	return fmt.Fprint(&g.buf, a...)
}

// Printf formats according to a format specifier for its operands and writes to the generated output.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(&g.buf, format, a...)
}

// Fprintln formats using the default formats to the generated output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(&g.buf, a...)
}

func (g *CodeGen) Reset() *CodeGen {
	g.buf.Reset()
	return g
}

func (g *CodeGen) Gen() *CodeGen {
	g.Println(`syntax = "proto3";`)
	g.Println()
	g.Printf("package %s;\n", g.PackageName)
	g.Println()

	if len(g.Options) > 0 {
		for k, v := range g.Options {
			g.Printf("option %s = \"%s\";\n", k, v)
		}
		g.Println()
	}

	if g.needGoogleProtobufTimestamp(g.Entity) {
		g.Println(`import "google/protobuf/timestamp.proto";`)
	}
	g.Println(`import "google/api/field_behavior.proto";`)
	g.Println(`import "protoc-gen-openapiv2/options/annotations.proto";`)
	g.Println()

	et := g.Entity
	structName := utils.PascalCase(et.Name)

	//* list
	g.Printf("message List%sRequest {\n", structName)
	g.genFields(et.Fields, nil, false)
	g.Println("}")
	g.Printf("message List%sReply {\n", structName)
	g.Println("  int64 total = 1 [(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }];")
	g.Println("  int64 page = 30 [(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }];")
	g.Println("  int64 per_page = 31 [(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }];")
	g.Printf("  repeated mapper.%s list = 32;\n", structName)
	g.Println("}")
	//* get
	g.Printf("message Get%sRequest {\n", structName)
	g.Println(`// @gotags: binding:"gt=0"`)
	g.Println("  int64 id = 1 [(google.api.field_behavior) = REQUIRED,")
	g.Println("  \t(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }];")
	g.Println("}")
	g.Printf("message Get%sReply {\n", structName)
	g.Printf("  mapper.%s %s = 1;\n", structName, utils.StyleName(g.Style, et.Name))
	g.Println("}")
	//* create
	g.Printf("message Add%sRequest {\n", structName)
	g.genFields(et.Fields, []string{"id", "created_at", "updated_at", "deleted_at"}, true)
	g.Println("}")
	g.Printf("message Add%sReply {\n", structName)
	g.Println("}")
	//* update
	g.Printf("message Update%sRequest {\n", structName)
	g.genFields(et.Fields, []string{"created_at", "updated_at", "deleted_at"}, true)
	g.Println("}")
	g.Printf("message Update%sReply {\n", structName)
	g.Println("}")
	//* delete
	g.Printf("message Delete%sRequest {\n", structName)
	g.Println(`// @gotags: binding:"gt=0"`)
	g.Println("  int64 id = 1 [(google.api.field_behavior) = REQUIRED,")
	g.Println("  \t(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }];")
	g.Println("}")
	g.Printf("message Delete%sReply {\n", structName)
	g.Println("}")
	//* bulk delete
	g.Printf("message BulkDelete%sRequest {\n", structName)
	g.Println(`// @gotags: binding:"required,dive,gt=0"`)
	g.Println("  repeated int64 id = 1 [(google.api.field_behavior) = REQUIRED];")
	g.Println("}")
	g.Printf("message BulkDelete%sReply {\n", structName)
	g.Println("}")

	return g
}

func (g *CodeGen) GenService() *CodeGen {
	et := g.Entity
	structName := utils.PascalCase(et.Name)

	g.Println(`syntax = "proto3";`)
	g.Println()
	g.Printf("package %s;\n", g.PackageName)
	g.Println()

	if len(g.Options) > 0 {
		for k, v := range g.Options {
			g.Printf("option %s = \"%s\";\n", k, v)
		}
		g.Println()
	}

	g.Println(`import "google/protobuf/empty.proto";`)
	g.Println(`import "google/api/annotations.proto";`)
	g.Println(`import "protoc-gen-openapiv2/options/annotations.proto";`)
	g.Println()

	g.Printf("service %s {\n", structName)
	g.Println("option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_tag) = {")
	g.Printf("description: \"%s\",\n", et.Comment)
	g.Println("};")

	//* list
	g.Printf("// 获取%s列表\n", et.Comment)
	g.Printf("rpc List%[1]s(List%[1]sRequest) returns (List%[1]sReply) {\n", structName)
	g.Println("option (google.api.http) = {")
	g.Printf("get: \"/v1/%s\"\n", et.Name)
	g.Println("};")
	g.Println("}")
	//* get
	g.Printf("// 获取%s\n", et.Comment)
	g.Printf("rpc Get%[1]s(Get%[1]sRequest) returns (Get%[1]sReply) {\n", structName)
	g.Println("option (google.api.http) = {")
	g.Printf("get: \"/v1/%s/{id}\"\n", et.Name)
	g.Println("};")
	g.Println("}")
	//* create
	g.Printf("// 增加%s\n", et.Comment)
	g.Printf("rpc Add%[1]s(Add%[1]sRequest) returns (Add%[1]sReply) {\n", structName)
	g.Println("option (google.api.http) = {")
	g.Printf("post: \"/v1/%s\"\n", et.Name)
	g.Println("body: \"*\"")
	g.Println("};")
	g.Println("}")
	//* update
	g.Printf("// 更新%s\n", et.Comment)
	g.Printf("rpc Update%[1]s(Update%[1]sRequest) returns (Update%[1]sReply) {\n", structName)
	g.Println("option (google.api.http) = {")
	g.Printf("put: \"/v1/%s\"\n", et.Name)
	g.Println("body: \"*\"")
	g.Println("};")
	g.Println("}")
	//* 删除
	g.Printf("// 删除%s\n", et.Comment)
	g.Printf("rpc Delete%[1]s(Delete%[1]sRequest) returns (Delete%[1]sReply) {\n", structName)
	g.Println("option (google.api.http) = {")
	g.Printf("delete: \"/v1/%s/{id}\"\n", et.Name)
	g.Println("};")
	g.Println("}")
	//* 批量删除
	g.Printf("// 批量删除%s\n", et.Comment)
	g.Printf("rpc BulkDelete%[1]s(BulkDelete%[1]sRequest) returns (BulkDelete%[1]sReply) {\n", structName)
	g.Println("option (google.api.http) = {")
	g.Printf("delete: \"/v1/%s\"\n", et.Name)
	g.Println("};")
	g.Println("}")
	g.Println("}")
	return g
}

func (g *CodeGen) genFields(fields []*proto.MessageField, skipName []string, required bool) {
	for i, m := range fields {
		if slices.Contains(skipName, m.ColumnName) {
			continue
		}
		if m.Comment != "" {
			g.Printf("  // %s\n", m.Comment)
		}
		annotations := []string{}
		if required {
			if g.isNumber(m.Type, m.TypeName) {
				g.Println(`// @gotags: binding:"gt=0"`)
			} else {
				g.Println(`// @gotags: binding:"required"`)
			}
			annotations = append(annotations, "(google.api.field_behavior) = REQUIRED")
		}
		typeName, tmpAnnotations := g.intoTypeNameAndAnnotation(m)
		annotations = append(annotations, tmpAnnotations...)
		annotation := ""
		if len(annotations) > 0 {
			annotation = fmt.Sprintf(" [%s]", strings.Join(annotations, ", "))
		}
		fieldName := utils.StyleName(g.Style, m.Name)

		seq := i + 1
		if m.Cardinality == protoreflect.Required {
			g.Printf("  %s %s = %d%s;\n", typeName, fieldName, seq, annotation)
		} else {
			g.Printf("  %s %s %s = %d%s;\n", m.Cardinality.String(), typeName, fieldName, seq, annotation)
		}
	}
}

func (g *CodeGen) intoTypeNameAndAnnotation(field *proto.MessageField) (string, []string) {
	annotations := make([]string, 0, 8)
	switch {
	case g.DisableBool && field.Type == protoreflect.BoolKind:
		return protoreflect.Int32Kind.String(), annotations
	case field.Type == protoreflect.MessageKind && field.TypeName == googleProtobufTimestamp:
		if g.DisableTimestamp {
			if g.EnableOpenapiv2Annotation {
				annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
			}
			return protoreflect.Int64Kind.String(), annotations
		} else {
			return field.TypeName, annotations
		}
	case (field.Type == protoreflect.Int64Kind || field.Type == protoreflect.Uint64Kind) && g.EnableOpenapiv2Annotation:
		annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
		fallthrough
	default:
		return field.Type.String(), annotations
	}
}

func (g *CodeGen) needGoogleProtobufTimestamp(message *proto.Message) bool {
	for _, v := range message.Fields {
		if !g.DisableTimestamp && v.Type == protoreflect.MessageKind && v.TypeName == googleProtobufTimestamp {
			return true
		}
	}
	return false
}

func (g *CodeGen) isNumber(typ protoreflect.Kind, typeName string) bool {
	switch typ {
	case protoreflect.Int32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Int64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Uint64Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.FloatKind,
		protoreflect.Sfixed64Kind,
		protoreflect.Fixed64Kind,
		protoreflect.DoubleKind:
		return true
	}
	return g.DisableTimestamp && typ == protoreflect.MessageKind && typeName == googleProtobufTimestamp
}
