package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"unicode"
)

type Field struct {
	Name string
	Type string
}

type Func struct {
	Recv    *Field
	Name    string
	Params  []Field
	Returns []Field
	Body    string
}

type GoFile struct {
	Package string
	Funcs   []*Func
}

func Parse(filename string, r io.Reader) (*GoFile, error) {
	fmt.Println(filename)
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
				Type: fmt.Sprintf("%s", field.Type),
			})
		}
	}
	return fieldSlice
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
