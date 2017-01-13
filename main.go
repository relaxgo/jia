// Package main provides ...
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"text/template"
	"unicode"
)

var (
	src string = `
package subtitle

import (
	"context"
	"github.com/relaxgo/subtitle-server/db"
)

// some des
func FindProjects(c context.Context, where, sorts string, skip int, limit int) ([]Project, error) {
	list := make([]Project, 0)

	query := db.DB.
		Preload("Ower").
		Offset(skip).
		Limit(limit).
		Find(&list)

	return list, query.Error
}

/*
* param
*/
func AddProject(c context.Context, project *Project) (*Project, error) {
	query := db.DB.Create(project)
	return project, query.Error
}
	`
)

var t *template.Template

func init() {
	tmpl := `
	func {{.Name }}(c echo.Context) error {
    ctx := ctx.New(c)
		{{ range $_, $p := .Params -}}
		{{ $p.Name }} := params.{{MethodOfType $p.Type }}(c, "{{$p.Name}}")
		{{ end }}
    v, err:= {{.Pkg}}.{{.Name}}(ctx, {{.ParamsText}})
    if err != nil {
        return err
    }
    return c.JSON(http.StatusOK, v)
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

func main() {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	bf := bytes.NewBuffer(nil)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			return true
		case *ast.FuncDecl:
			if data, err := GenRouter(f, x); err != nil {
				fmt.Println(err)
			} else {
				bf.Write(data)
			}
		}
		return false
	})
	data, err := format.Source(bf.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
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
	list := fn.Type.Params.List[1:]
	var params []Param
	for _, field := range list {
		t := fmt.Sprintf("%s", field.Type)
		if t != "int" && t != "string" {
			t = "object"
		}
		for _, n := range field.Names {
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
