package {{.Package}}

{{- $e := .Entity}}
{{- $stName := pascalcase $e.Name}}

type Update{{$stName}}ByPartial struct {
{{- range $f := $e.Fields}}
    {{if or (eq $f.GoName "CreatedAt") (eq $f.GoName "UpdatedAt") (eq $f.GoName "DeletedAt")}}// {{end}}{{$f.GoName}} {{if ne $f.GoName "Id"}}*{{- end}}{{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
{{- end}}
}

type Get{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{if or (eq $f.GoName "CreatedAt") (eq $f.GoName "UpdatedAt") (eq $f.GoName "DeletedAt")}}// {{end}}{{$f.GoName}} {{if eq $f.Type.Type 1 }}*{{- end}}{{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
{{- end}}
}

type Exist{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{if or (eq $f.GoName "CreatedAt") (eq $f.GoName "UpdatedAt") (eq $f.GoName "DeletedAt")}}// {{end}}{{$f.GoName}} {{if eq $f.Type.Type 1 }}*{{- end}}{{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
{{- end}}
}

type List{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{if or (eq $f.GoName "CreatedAt") (eq $f.GoName "UpdatedAt") (eq $f.GoName "DeletedAt")}}// {{end}}{{$f.GoName}} {{if eq $f.Type.Type 1 }}*{{- end}}{{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
{{- end}}
    Page    int64 `json:"page"`
    PerPage int64 `json:"perPage"`
}

type PluckId{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{if or (eq $f.GoName "CreatedAt") (eq $f.GoName "UpdatedAt") (eq $f.GoName "DeletedAt")}}// {{end}}{{$f.GoName}} {{if eq $f.Type.Type 1 }}*{{- end}}{{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
{{- end}}
}