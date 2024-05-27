package {{.Package}}

{{- $e := .Entity}}
{{- $stName := pascalcase $e.Name}}

type Update{{$stName}}ByPartial struct {
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
    {{$f.GoName}} {{if ne $f.GoName "Id"}}*{{- end}}{{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
    {{- end}}
{{- end}}
}

type Get{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
    {{$f.GoName}} {{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
    {{- end}}
{{- end}}
}

type Exist{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
    {{$f.GoName}} {{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
    {{- end}}
{{- end}}
}

type List{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
    {{$f.GoName}} {{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
    {{- end}}
{{- end}}
    Page    int64 `json:"page"`
    PerPage int64 `json:"perPage"`
}

type PluckId{{$stName}}ByFilter struct {
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
    {{$f.GoName}} {{$f.Type.Ident}} `json:"{{smallcamelcase $f.ColumnName}}"`
    {{- end}}
{{- end}}
}