package jia

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"unicode"

	"gopkg.in/yaml.v1"
)

type FieldType types.Type

type Field struct {
	Name string
	Type FieldType
}

type Struct struct {
	Name   string
	Fields []Field
}

type Func struct {
	Recv      *Field
	Name      string
	Params    []Field
	Returns   []Field
	Body      string
	Doc       string
	ParsedDoc map[string]interface{}
}

type GoFile struct {
	Package string
	Funcs   []*Func
	Structs []Struct
}

func (f *Field) TypeKind() string {
	s := fmt.Sprintf("%T", f.Type)
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return s[i+1:]
		}
	}
	return ""
}

func (f *Field) TypeName() string {
	return fmt.Sprint(f.Type)
}

func (f *Field) IsBasic() bool {
	return f.TypeKind() == "Basic"
}

func ParsePackage(filename string) (*GoFile, error) {
	fset := token.NewFileSet()
	dir, pkg, _ := splitPath(filename)
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	targetPkg := pkgs[pkg]
	if targetPkg == nil {
		return nil, errors.New("not have package: " + pkg)
	}
	return Parse(filename, fset, targetPkg.Files)
}

func ParseFile(filename string, r io.Reader) (gofile *GoFile, err error) {
	if r == nil {
		r, err = os.Open(filename)
		if err != nil {
			return nil, errors.New("parser.OpenFile:" + err.Error())
		}
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, r, parser.ParseComments)
	if err != nil {
		return nil, errors.New("parser.ParseFile:" + err.Error())
	}
	return Parse(filename, fset, map[string]*ast.File{filename: f})
}

func Parse(filename string, fset *token.FileSet, filemap map[string]*ast.File) (*GoFile, error) {
	conf := types.Config{Importer: importer.Default()}
	// TODO import from source ?
	//imp := NewImporter(&build.Default, fset, make(map[string]*types.Package))
	//conf := types.Config{Importer: imp}
	info := &types.Info{
		Defs:  make(map[*ast.Ident]types.Object),
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	_, pkg, _ := splitPath(filename)
	f := filemap[filename]
	if f == nil {
		return nil, errors.New("file not exist")
	}

	_, err := conf.Check(pkg, fset, mapToSlice(filemap), info)
	if err != nil {
		// FIXME Check can't Resolve type in self pk
		log.Println(errors.New("conf.Check:" + err.Error()))
	}

	file := &GoFile{}
	file.Package = f.Name.String()

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			return true
		case *ast.FuncDecl:
			file.Funcs = append(file.Funcs, ParseFunc(info, x))
		}
		return false
	})
	return file, nil
}

func ParseFunc(info *types.Info, funcDecl *ast.FuncDecl) *Func {
	f := &Func{}

	if funcDecl.Recv != nil {
		f.Recv = &ParseFields(info, funcDecl.Recv)[0]
	}

	f.Name = funcDecl.Name.Name
	f.Params = ParseFields(info, funcDecl.Type.Params)
	f.Returns = ParseFields(info, funcDecl.Type.Results)
	// TODO doc format should be set in cmd
	if funcDecl.Doc != nil {
		f.Doc = funcDecl.Doc.Text()
		f.ParsedDoc = make(map[string]interface{})
		yaml.Unmarshal([]byte(funcDecl.Doc.Text()), &f.ParsedDoc)
	}

	return f
}

func ParseFields(info *types.Info, list *ast.FieldList) []Field {
	var fieldSlice []Field
	if list == nil {
		return nil
	}

	for _, field := range list.List {
		for _, n := range field.Names {
			fieldSlice = append(fieldSlice, Field{
				Name: n.Name,
				Type: ParseFieldType(info, field.Type),
			})
		}
	}

	return fieldSlice
}

func Underlying(t types.Type) {
	switch s := t.(type) {
	case *types.Named:
		fmt.Println(s.Obj().Name())
		// Underlying(s.Underlying())
	case *types.Struct:
		// fmt.Println(s)
		for i := 0; i < s.NumFields(); i++ {
			// s := s.Field(i)
			// Underlying(s.Type())
		}
	case *types.Slice:
		// Underlying(s.Elem())
	}
}

func ParseFieldType(info *types.Info, n ast.Expr) FieldType {
	ts := info.Types
	if t, ok := ts[n]; ok {
		return FieldType(t.Type)
	}
	return nil
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

func splitPath(f string) (dir, pkg, gofile string) {
	dir, gofile = path.Split(f)
	dirs := strings.Split(strings.TrimRight(dir, "/"), "/")

	if l := len(dirs); l != 0 {
		pkg = dirs[len(dirs)-1]
	}
	return
}

func mapToSlice(m map[string]*ast.File) []*ast.File {
	list := make([]*ast.File, 0)
	for _, v := range m {
		if v != nil {
			list = append(list, v)
		}
	}
	return list
}
