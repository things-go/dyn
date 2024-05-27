package command

import (
	"embed"
	"errors"
	"text/template"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

//go:embed template/*.tpl
var Static embed.FS

var TemplateFuncs = template.FuncMap{
	"add":            func(a, b int) int { return a + b },
	"snakecase":      func(s string) string { return utils.SnakeCase(s) },
	"kebabcase":      func(s string) string { return utils.Kebab(s) },
	"pascalcase":     func(s string) string { return utils.PascalCase(s) },
	"smallcamelcase": func(s string) string { return utils.SmallCamelCase(s) },
}
var (
	tpl = template.Must(template.New("components").
		Funcs(TemplateFuncs).
		ParseFS(Static, "template/*.tpl"))
	dalRapierTpl = tpl.Lookup("dal_rapier.tpl")
	dalGormTpl   = tpl.Lookup("dal_gorm.tpl")
	dalQueryTpl  = tpl.Lookup("dal_query.tpl")
	dalOptionTpl = tpl.Lookup("dal_option.tpl")
)

type Dal struct {
	Package     string
	Imports     []string
	ModelPrefix string
	QueryPrefix string
	RepoPrefix  string
	Entity      *ens.EntityDescriptor
}

type DalQuery struct {
	PackageName    string
	Imports        []string
	ModelQualifier string
	Entity         ens.EntityDescriptor
}

func GetUsedTemplate(t string) (*template.Template, error) {
	switch t {
	case "builtin-gorm":
		return dalGormTpl, nil
	case "builtin-rapier":
		return dalRapierTpl, nil
	default:
		t, err := ParseTemplateFromFile(t)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
}

func ParseTemplateFromFile(filename string) (*template.Template, error) {
	if filename == "" {
		return nil, errors.New("required template filename")
	}
	tt, err := template.New("custom").
		Funcs(TemplateFuncs).
		ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	ts := tt.Templates()
	if len(ts) == 0 {
		return nil, errors.New("not found any template")
	}
	return ts[0], nil
}
