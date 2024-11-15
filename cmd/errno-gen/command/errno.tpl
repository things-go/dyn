// Code generated by errno-gen. DO NOT EDIT.
// version: {{.Version}}
{{- if .IsDeprecated}}
//
// Deprecated: this is a deprecated file.
{{- end}}
package {{.Package}}

import (
	errors "{{.Epk}}"
)
{{- range $e := .Enums}}

{{- range $ee := $e.Values}}
// Err{{$ee.OriginalName}} {{$ee.Value}}: {{.Mapping}}
func Err{{$ee.OriginalName}}(opts ...errors.Option) *errors.Error {
	return errors.New(int32({{$ee.OriginalName}}), {{$ee.OriginalName}}.String(), opts...)
}
{{- end}}

{{- end}}