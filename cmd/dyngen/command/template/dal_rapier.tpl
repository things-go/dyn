package {{.Package}}

import (
    "context"
    
    "gorm.io/gorm"
    rapier "github.com/thinkgos/gorm-rapier"

{{- range $e := .Imports}}
    {{- if ne $e ""}}
    "{{$e}}"
    {{- end}}
{{- end}}
)

{{- $e := .Entity}}
{{- $stName := pascalcase $e.Name}}
{{- $mdPrefix := .ModelPrefix}}
{{- $queryPrefix := .QueryPrefix}}
{{- $repoPrefix := .RepoPrefix}}
{{- $mdName := printf "%s%s" $mdPrefix $stName}}

var _ {{$stName}}Dal = {{$stName}}{}
var _ = rapier.NewExecutor[{{$stName}}]

type {{$stName}}Dal interface {
    Create(ctx context.Context, v ...*{{$mdName}}) (int64, error)
    Delete(ctx context.Context, id ...int64) (int64, error)
    UpdateFull(ctx context.Context, v *{{$mdName}}) (int64, error)
    UpdatePartial(ctx context.Context, v *{{$queryPrefix}}Update{{$stName}}ByPartial) (int64, error) 
    Get(ctx context.Context, id int64, funcs ...DalCondition) (*{{$mdName}}, error)
    GetByFilter(ctx context.Context, q *{{$queryPrefix}}Get{{$stName}}ByFilter, funcs ...DalCondition) (*{{$mdName}}, error) 
    ExistByFilter(ctx context.Context, q *{{$queryPrefix}}Exist{{$stName}}ByFilter, funcs ...DalCondition) (bool, error) 
    Count(ctx context.Context, q *{{$queryPrefix}}List{{$stName}}ByFilter) (int64, error) 
    List(ctx context.Context, q *{{$queryPrefix}}List{{$stName}}ByFilter) ([]*{{$mdName}}, error)
    ListPage(ctx context.Context, q *{{$queryPrefix}}List{{$stName}}ByFilter) ([]*{{$mdName}}, int64, error)
    PluckIdByFilter(ctx context.Context, q *{{$queryPrefix}}PluckId{{$stName}}ByFilter) ([]int64, error) 
}

type {{$stName}} struct {
	db *gorm.DB
}

func New{{$stName}}(db *gorm.DB) {{$stName}}Dal {
    return {{$stName}} {
        db: db,
    }
}

func (b {{$stName}}) Create(ctx context.Context, v ...*{{$mdName}}) (int64, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Create(v...)
}

func (b {{$stName}}) Delete(ctx context.Context, id ...int64) (int64, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            Where(ref.Id.In(id...)).
            Delete()
}


func (b {{$stName}}) UpdateFull(ctx context.Context, v *{{$mdName}}) (int64, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            Select("*").
            Where(ref.Id.Eq(v.Id)).
            Updates(v)
}

func (b {{$stName}}) UpdatePartial(ctx context.Context, v *{{$queryPrefix}}Update{{$stName}}ByPartial) (int64, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
	ct := defaultPool.Get()
	defer defaultPool.Put(ct)
	up := ct.Exprs
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt") (ne $f.GoName "Id")}}
    if v.{{$f.GoName}} != nil {
        up = append(up, ref.{{$f.GoName}}.Value(*v.{{$f.GoName}}))
    }
    {{- end}}
{{- end}}
    if len(up) == 0 {
        return 0, nil
    }
    return ref.New_Executor(b.db).Model().
            Where(ref.Id.Eq(v.Id)).
            UpdatesExpr(up...)
}

func (b {{$stName}}) Get(ctx context.Context, id int64, funcs ...DalCondition) (*{{$mdName}}, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            SelectExpr(ref.Select_Expr()...).
            Scopes(funcs...).
            Where(ref.Id.Eq(id)).
            TakeOne()
}

func (b {{$stName}}) GetByFilter(ctx context.Context, q *{{$queryPrefix}}Get{{$stName}}ByFilter, funcs ...DalCondition) (*{{$mdName}}, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            SelectExpr(ref.Select_Expr()...).
            Scopes(funcs...).
            Scopes(get{{$stName}}Filter(ref, q)).
            TakeOne()
}

func (b {{$stName}}) ExistByFilter(ctx context.Context, q *{{$queryPrefix}}Exist{{$stName}}ByFilter, funcs ...DalCondition) (existed bool, err error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            Scopes(funcs...).
            Scopes(func(db *gorm.DB) *gorm.DB {
        {{- range $f := $e.Fields}}
            {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
            {{- if eq $f.Type.Type 15 }}
                if q.{{$f.GoName}} != "" {
            {{- else if eq $f.Type.Type 18 }}
                if !q.{{$f.GoName}}.IsZero() {
            {{- else if eq $f.Type.Type 1 }}
                if q.{{$f.GoName}} != nil {
            {{- else }}
                if q.{{$f.GoName}} != 0 {
            {{- end}}
                    db = db.Where(ref.{{$f.GoName}}.Eq({{if eq $f.Type.Type 1 }}*{{- end}}q.{{$f.GoName}}))
                }
            {{- end}}
        {{- end}}
                return db
            }).
            Exist()
}

func (b {{$stName}}) Count(ctx context.Context, q *{{$queryPrefix}}List{{$stName}}ByFilter) (total int64, err error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            Scopes(list{{$stName}}Filter(ref, q)).
            Count()
}


func (b {{$stName}}) List(ctx context.Context, q *{{$queryPrefix}}List{{$stName}}ByFilter) ([]*{{$mdName}}, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            SelectExpr(ref.Select_Expr()...).
            Scopes(list{{$stName}}Filter(ref, q), Limit(q.Page, q.PerPage)).
            FindAll()
}

func (b {{$stName}}) ListPage(ctx context.Context, q *{{$queryPrefix}}List{{$stName}}ByFilter) ([]*{{$mdName}}, int64, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
	return ref.New_Executor(b.db).Model().
        SelectExpr(ref.Select_Expr()...).
		Scopes(list{{$stName}}Filter(ref, q)).
		FindAllPaginate(q.Page, q.PerPage)
}

func (b {{$stName}}) PluckIdByFilter(ctx context.Context, q *{{$queryPrefix}}PluckId{{$stName}}ByFilter) ([]int64, error) {
    ref := {{$repoPrefix}}Ref_{{$stName}}()
    return ref.New_Executor(b.db).Model().
            Scopes(func(db *gorm.DB) *gorm.DB {
        {{- range $f := $e.Fields}}
            {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
            {{- if eq $f.Type.Type 15 }}
                if q.{{$f.GoName}} != "" {
            {{- else if eq $f.Type.Type 18 }}
                if !q.{{$f.GoName}}.IsZero() {
            {{- else if eq $f.Type.Type 1 }}
                if q.{{$f.GoName}} != nil {
            {{- else }}
                if q.{{$f.GoName}} != 0 {
            {{- end}}
                    db = db.Where(ref.{{$f.GoName}}.Eq({{if eq $f.Type.Type 1 }}*{{- end}}q.{{$f.GoName}}))
                }
            {{- end}}
        {{- end}}
                return db
            }).
            PluckExprInt64(ref.Id)
}

func get{{$stName}}Filter(ref *{{$repoPrefix}}{{$stName}}_Native, q *{{$queryPrefix}}Get{{$stName}}ByFilter) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        {{- range $f := $e.Fields}}
            {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
            {{- if eq $f.Type.Type 15 }}
                if q.{{$f.GoName}} != "" {
            {{- else if eq $f.Type.Type 18 }}
                if !q.{{$f.GoName}}.IsZero() {
            {{- else if eq $f.Type.Type 1 }}
                if q.{{$f.GoName}} != nil {
            {{- else }}
                if q.{{$f.GoName}} != 0 {
            {{- end}}
                    db = db.Where(ref.{{$f.GoName}}.Eq({{if eq $f.Type.Type 1 }}*{{- end}}q.{{$f.GoName}}))
                }
            {{- end}}
        {{- end}}
                return db
            }
}

func list{{$stName}}Filter(ref *{{$repoPrefix}}{{$stName}}_Native, q *{{$queryPrefix}}List{{$stName}}ByFilter) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
{{- range $f := $e.Fields}}
    {{- if and (ne $f.GoName "CreatedAt") (ne $f.GoName "UpdatedAt") (ne $f.GoName "DeletedAt")}}
    {{- if eq $f.Type.Type 15 }}
        if q.{{$f.GoName}} != "" {
    {{- else if eq $f.Type.Type 18 }}
        if !q.{{$f.GoName}}.IsZero() {
    {{- else if eq $f.Type.Type 1 }}
        if q.{{$f.GoName}} != nil {
    {{- else }}
        if q.{{$f.GoName}} != 0 {
    {{- end}}
            db = db.Where(ref.{{$f.GoName}}.Eq({{if eq $f.Type.Type 1 }}*{{- end}}q.{{$f.GoName}}))
        }
    {{- end}}
{{- end}}
        return db
    }
}