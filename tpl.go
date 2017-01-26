package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"reflect"
	"strings"
	"unicode"
)

var EchoTemp *template.Template

func init() {
	tmpl := `
	package controls

	import (
		"bitbucket.org/reewoow_web/jiyin/models"
		"bitbucket.org/reewoow_web/jiyin/components/param"
		"github.com/labstack/echo"
	)
	{{ $funcs := .ValidFuncs }}
  {{ range $_, $f := $funcs }}
	func {{$f.Name }}(c echo.Context) error {
		{{ range $_, $p := $f.Params -}}
		{{ $p.Name }} := param.{{MethodOfType $p.Type }}(c, "{{ ToLower $p.Name}}")
		{{ end }}
    v, err:= {{$.Package}}.{{$f.Name}}({{JoinFiled $f.Params "Name" ","}})
		return Respond(c, v, err)
	}
	{{ end }}
	`

	t := template.New("fn").
		Funcs(template.FuncMap{
		"MethodOfType": func(str string) string {
			for i, v := range str {
				return string(unicode.ToUpper(v)) + str[i+1:]
			}
			return ""
		},
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
		"JoinFiled": func(slice interface{}, fieldName, sep string) string {
			v := reflect.ValueOf(slice)
			t := v.Type()
			if t.Kind() != reflect.Slice {
				panic("JoinFiled need slice")
			}
			l := v.Len()
			s := make([]string, l, l)
			for i := 0; i < l; i++ {
				s[i] = v.Index(i).FieldByName(fieldName).String()
			}
			return strings.Join(s, sep)
		},
	})

	_, err := t.Parse(tmpl)
	fmt.Println(err)
	EchoTemp = t
}

func Render(f *GoFile, t *template.Template) ([]byte, error) {
	bf := bytes.NewBuffer(nil)
	err := EchoTemp.Execute(bf, f)

	if err != nil {
		return nil, err
	}
	return format.Source(bf.Bytes())
}
