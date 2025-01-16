package protoenum

import (
	"github.com/things-go/proc/proc"
	"google.golang.org/protobuf/compiler/protogen"
)

// annotation const value
const (
	Identity               = "enum"
	Attribute_Name_Mapping = "mapping"
	Attribute_Name_Label   = "label"
)

type EnumDerive struct {
	Enabled bool
}

func ParseDeriveEnum(s protogen.Comments) (*EnumDerive, proc.CommentLines) {
	ret := &EnumDerive{Enabled: false}
	derives, remainComments := proc.NewCommentLines(string(s)).FindDerives(Identity)
	ret.Enabled = proc.Derives(derives).ContainHeadless(Identity)
	return ret, remainComments
}

type EnumValueDerive struct {
	Mapping string
	Label   string
}

func ParseDeriveEnumValue(s protogen.Comments) (*EnumValueDerive, proc.CommentLines) {
	ret := &EnumValueDerive{Mapping: "", Label: ""}
	derives, remainComments := proc.NewCommentLines(string(s)).FindDerives(Identity)
	mappingValues := proc.Derives(derives).FindValue(Identity, Attribute_Name_Mapping)
	for _, v := range mappingValues {
		if v, ok := v.(proc.String); ok {
			ret.Mapping = v.Value
			break
		}
	}
	labelValues := proc.Derives(derives).FindValue(Identity, Attribute_Name_Label)
	for _, v := range labelValues {
		if v, ok := v.(proc.String); ok {
			ret.Label = v.Value
			break
		}
	}
	return ret, remainComments
}
