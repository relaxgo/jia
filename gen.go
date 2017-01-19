package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"text/template"
	"unicode"
)

var t *template.Template

var (
	header = []byte(`
package controls

import (
	"bitbucket.org/reewoow_web/jiyin/models"
	"bitbucket.org/reewoow_web/jiyin/components/param"
	"github.com/labstack/echo"
)
	`)
)

func init() {
	tmpl := `
	func {{.Name }}(c echo.Context) error {
		{{ range $_, $p := .Params -}}
		{{ $p.Name }} := param.{{MethodOfType $p.Type }}(c, "{{$p.Name}}")
		{{ end }}
    v, err:= {{.Pkg}}.{{.Name}}({{.ParamsText}})
		return Respond(c, v, err)
	}
	`
	t = template.New("fn").
		Funcs(template.FuncMap{
		"MethodOfType": func(name string) string {
			return UpperInitial(name)
		},
	})
	_, err := t.Parse(tmpl)
	fmt.Println(err)
}

func Gen(src io.Reader) []byte {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	bf := bytes.NewBuffer(header)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			return true
		case *ast.FuncDecl:
			// TODO 过滤特定函数名, 如小写开头，或者ignore
			if unicode.IsLower(rune(x.Name.Name[0])) {
				return false
			}
			if data, err := GenRouter(f, x); err != nil {
				fmt.Println(err)
			} else {
				bf.Write(data)
			}
		}
		return false
	})
	data, err := format.Source(bf.Bytes())
	handleErr(err)
	return data
}

type RouteData struct {
	Pkg    string
	Name   string
	Params []Param
}

type Param struct {
	Name string
	Type string
}

func GenRouter(f *ast.File, funcDecl *ast.FuncDecl) ([]byte, error) {
	r := NewRoute(f, funcDecl)
	bf := bytes.NewBuffer(nil)
	err := t.Execute(bf, r)
	return bf.Bytes(), err
}

func NewRoute(file *ast.File, fn *ast.FuncDecl) *RouteData {
	list := fn.Type.Params.List[0:]
	var params []Param
	for _, field := range list {
		t := fmt.Sprintf("%s", field.Type)
		if t != "int" && t != "string" {
			t = "object"
		}
		for _, n := range field.Names {
			// TODO 自定义参数名
			params = append(params, Param{
				Name: n.Name,
				Type: t,
			})
		}
	}
	return &RouteData{
		Pkg:    file.Name.String(),
		Name:   fn.Name.Name,
		Params: params,
	}
}

func (r *RouteData) ParamsText() string {
	if r.Params == nil {
		return ""
	}
	strs := make([]string, len(r.Params))
	for i, p := range r.Params {
		strs[i] = p.Name
	}
	return strings.Join(strs, ",")
}

func UpperInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
