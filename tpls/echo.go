package tpls

import (
	"fmt"
	"html/template"

	"github.com/relaxgo/jia"
)

var (
	EchoTemp = template.New("root").Funcs(jia.BaseFuncs)
)

func init() {
	tmpl := `
	package controls

	import (
		"net/http"

		"github.com/relaxgo/tangram/param"
		"github.com/labstack/echo"
	)

	{{ $funcs := .ValidFuncs }}
  {{ range $_, $f := $funcs }}
	func {{$f.Name }}(c echo.Context) error {
		{{ range $_, $p := $f.Params -}}
		{{ $p.Name }} := param.{{firstToUpper $p.Type }}(c, "{{ toLower $p.Name}}")
		{{ end }}
    v, err:= {{$.Package}}.{{$f.Name}}({{joinField $f.Params "Name" ","}})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, v)
	}
	{{ end }}
	`
	_, err := EchoTemp.Parse(tmpl)
	if err != nil {
		fmt.Println(err)
	}
}
