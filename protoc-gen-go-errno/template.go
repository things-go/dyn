package main

import (
	"bytes"
	"text/template"
)

var errnoTemplate = `
{{ range .Errors }}
func Is{{.CamelValue}}(err error) bool {
	e := errors.FromError(err)
{{- if or (eq .Code 400) (eq .Code 500)}}
	return e.Code == {{.Code}}
{{- else}}
	return e.Detail == {{.Name}}_{{.Value}}.String() && e.Code == {{.Code}} 
{{- end}}
}
func Err{{.CamelValue}}({{if or (eq .Code 400) (eq .Code 500)}}detail string{{else}}message ...string{{end}}) *errors.Error {
{{- if or (eq .Code 400) (eq .Code 500)}}
	return errors.New({{.Code}}, "{{.Message}}", detail)
{{- else}}
	s := "{{.Message}}"
	if len(message) > 0 {
		s = message[0] 
	}
	return errors.New({{.Code}}, s, {{.Name}}_{{.Value}}.String())
{{- end}}
}
func Err{{.CamelValue}}f(format string, args ...interface{}) *errors.Error {
{{- if or (eq .Code 400) (eq .Code 500)}}
	 return errors.New({{.Code}}, "{{.Message}}", fmt.Sprintf(format, args...))
{{- else}}
	 return errors.New({{.Code}}, fmt.Sprintf(format, args...), {{.Name}}_{{.Value}}.String())
{{- end}}
}
{{- end }}
`

type errorInfo struct {
	Name       string
	Code       int
	Value      string
	CamelValue string
	Message    string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errno").Parse(errnoTemplate)
	if err != nil {
		panic(err)
	}
	if err = tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
