package command

import (
	"embed"
	"html/template"
	"io"
)

//go:embed errno.tpl
var Static embed.FS

var errorsTemplate = template.Must(template.New("components").
	ParseFS(Static, "errno.tpl")).
	Lookup("errno.tpl")

type GenFile struct {
	Version      string
	IsDeprecated bool
	Package      string
	Epk          string
	Enums        []*Enumerate
}

type Enumerate struct {
	Type     string
	TypeName string
	Explain  string
	Values   []*Value
}

func (e *GenFile) execute(w io.Writer) error {
	return errorsTemplate.Execute(w, e)
}
