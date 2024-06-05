package command

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type Package struct {
	Name  string
	Defs  map[*ast.Ident]types.Object
	Files []*File
}

type File struct {
	Pkg         *Package
	File        *ast.File
	TypeName    string
	TypeComment string
	Type        string
	Values      []*Value
}

// Value represents a declared constant.
type Value struct {
	OriginalName string // 常量定义的名称
	Mapping      string // 注释名称, 如果没有, 则同常量名称
	// value相关
	Value  uint64 // 需要时转为`int64`.
	Signed bool   // `constant`是否是有符号类型.
	Val    string // `constant`的字符串值,由"go/constant"包提供.
}

func (v *Value) String() string { return v.Val }

// SortValue 使我们可以将`constants`进行排序, 要谨慎地按照有符号或无符号的顺序进行恰当的处理
type SortValue []*Value

func (b SortValue) Len() int      { return len(b) }
func (b SortValue) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b SortValue) Less(i, j int) bool {
	if b[i].Signed {
		return int64(b[i].Value) < int64(b[j].Value)
	}
	return b[i].Value < b[j].Value
}

// GenDecl processes one declaration clause.
func (f *File) GenDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || !slices.Contains([]token.Token{token.CONST, token.TYPE}, decl.Tok) {
		// 只关心 const, Type 声明
		return true
	}

	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		if decl.Tok == token.TYPE { // type spec
			tsepc := spec.(*ast.TypeSpec)          // 必定是 TYPE
			if tsepc.Name.String() == f.TypeName { // 找到这个类型
				obj, ok := f.Pkg.Defs[tsepc.Name]
				if !ok {
					slog.Error(fmt.Sprintf("no value for type %s", f.TypeName))
					os.Exit(1)
				}
				basic := obj.Type().Underlying().(*types.Basic)
				if basic.Info()&types.IsInteger == 0 {
					slog.Error(fmt.Sprintf("can't handle non-integer constant type %s", typ))
					os.Exit(1)
				}
				f.Type = basic.Name() // 拿到类型
				if c := tsepc.Comment.Text(); c != "" {
					f.TypeComment += strings.TrimSuffix(strings.TrimSpace(c), "\n")
				}
			}
		} else { // const spec
			vspec := spec.(*ast.ValueSpec) // 必定是 CONST.
			if vspec.Type == nil && len(vspec.Values) > 0 {
				// "X = 1". With no type but a value. If the constant is untyped,
				// skip this vspec and reset the remembered type.
				typ = ""

				// If this is a simple type conversion, remember the type.
				// We don't mind if this is actually a call; a qualified call won't
				// be matched (that will be SelectorExpr, not Ident), and only unusual
				// situations will result in a function call that appears to be
				// a type conversion.
				ce, ok := vspec.Values[0].(*ast.CallExpr)
				if !ok {
					continue
				}
				id, ok := ce.Fun.(*ast.Ident)
				if !ok {
					continue
				}
				typ = id.Name
			}
			if vspec.Type != nil {
				// "X T". We have a type. Remember it.
				ident, ok := vspec.Type.(*ast.Ident)
				if !ok {
					continue
				}
				typ = ident.Name
			}
			if typ != f.TypeName {
				// This is not the type we're looking for.
				continue
			}
			// We now have a list of names (from one line of source code) all being
			// declared with the desired type.
			// Grab their names and actual values and store them in f.values.
			for _, name := range vspec.Names {
				if name.Name == "_" {
					continue
				}
				// This dance lets the type checker find the values for us. It's a
				// bit tricky: look up the object declared by the name, find its
				// types.Const, and extract its value.
				obj, ok := f.Pkg.Defs[name]
				if !ok {
					slog.Error(fmt.Sprintf("no value for constant %s", name))
					os.Exit(1)
				}
				info := obj.Type().Underlying().(*types.Basic).Info()
				if info&types.IsInteger == 0 {
					slog.Error(fmt.Sprintf("can't handle non-integer constant type %s", typ))
					os.Exit(1)
				}
				value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
				if value.Kind() != constant.Int {
					slog.Error(fmt.Sprintf("can't happen: constant is not an integer %s", name))
					os.Exit(1)
				}
				i64, isInt := constant.Int64Val(value)
				u64, isUint := constant.Uint64Val(value)
				if !isInt && !isUint {
					slog.Error(fmt.Sprintf("internal error: value of %s is not an integer: %s", name, value.String()))
					os.Exit(1)
				}
				if !isInt {
					u64 = uint64(i64)
				}
				v := &Value{
					OriginalName: name.Name,
					Value:        u64,
					Signed:       info&types.IsUnsigned == 0,
					Val:          value.String(),
				}
				if c := vspec.Comment; c != nil && len(c.List) == 1 {
					v.Mapping = strings.TrimSpace(c.Text())
				} else {
					v.Mapping = v.OriginalName
				}
				f.Values = append(f.Values, v)
			}
		}
	}
	return false
}
