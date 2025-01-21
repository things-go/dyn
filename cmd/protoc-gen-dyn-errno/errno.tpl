{{ range .Errors }}
func Is{{.CamelValue}}(err error) bool {
	return errorx.EqualCode(err, {{.Code}})
}
func Err{{.CamelValue}}() *errorx.Error {
	return errorx.New({{.Code}}, "{{.Message}}")
}
{{- end }}