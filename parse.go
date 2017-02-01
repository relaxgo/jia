package jia

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"unicode"

	"gopkg.in/yaml.v2"
)

type Field struct {
	Name string
	Type FieldType
}

type Func struct {
	Recv    *Field
	Name    string
	Params  []Field
	Returns []Field
	Body    string
	Doc     map[string]interface{}
}

type GoFile struct {
	Package string
	Funcs   []*Func
}

func Parse(filename string, r io.Reader) (*GoFile, error) {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, filename, r, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	file := &GoFile{}
	file.Package = f.Name.String()
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			return true
		case *ast.FuncDecl:
			file.Funcs = append(file.Funcs, ParseFunc(x))
		}
		return false
	})
	return file, nil
}

func ParseFunc(funcDecl *ast.FuncDecl) *Func {
	f := &Func{}

	if funcDecl.Recv != nil {
		f.Recv = &ParseFields(funcDecl.Recv)[0]
	}

	f.Name = funcDecl.Name.Name
	f.Params = ParseFields(funcDecl.Type.Params)
	f.Returns = ParseFields(funcDecl.Type.Results)

	if funcDecl.Doc != nil {
		yaml.Unmarshal([]byte(funcDecl.Doc.Text()), &f.Doc)
	}

	return f
}

func ParseFields(list *ast.FieldList) []Field {
	var fieldSlice []Field
	if list == nil {
		return nil
	}
	for _, field := range list.List {
		for _, n := range field.Names {

			fieldSlice = append(fieldSlice, Field{
				Name: n.Name,
				Type: *ParseField(&FieldType{}, field.Type),
			})
		}
	}
	return fieldSlice
}

type (
	FieldType struct {
		Base  bool
		Point bool
		Pkg   string
		Name  string
	}
)

func ParseField(t *FieldType, n ast.Expr) *FieldType {
	switch e := n.(type) {
	case *ast.Ident:
		t.Base = true
		t.Name = fmt.Sprintf("%s", e)
	case *ast.StarExpr:
		ParseField(t, e.X)
		t.Point = true
		t.Base = false
	case *ast.SelectorExpr:
		ParseField(t, e.Sel)
		t.Pkg = fmt.Sprintf("%s", e)
		t.Base = false
	default:
		panic(fmt.Sprintf("Unsupport Expr: %T", e))
	}
	return t
}

func (file *GoFile) ValidFuncs() []Func {
	var fns []Func
	for _, f := range file.Funcs {
		if f.Recv != nil {
			continue
		}
		if unicode.IsLower(rune(f.Name[0])) {
			continue
		}
		fns = append(fns, *f)
	}
	return fns
}
