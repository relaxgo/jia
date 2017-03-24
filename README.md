## Jia
Parse the go file by libraries like `go/ast`, then generate the other file that related to this go file by your custom template.

[中文文档](README-zh_CN.md)

## Warning
It can be used now, but still in progress. All API may be changed.

## Example
Write you web api as a function. Pass all param directly.
```go
package models

type Order struct {
	Id    int
	Title string
	Price float64
}

func FindOrderById(userid, orderid int) (*Order, error) {
	// ...
	return &Order{}, nil
}

func CreateOrder(userid int, title string, price float64) (*Order, error) {
	// ...
	return &Order{}, nil
}
```
Add you custom route template.
```html
// generate by jia
package routes

import (
  "net/http"

  "github.com/relaxgo/tangram/param"
  "somepkg/models"
)

{{ range $_, $f := .ValidFuncs }}
{{if not $f.ParsedDoc.api_ignore }}
func {{ $f.Name }}(w http.ResponseWriter, r *http.Request) {
  p := NewRequestValue(r)

  {{ range $_, $p := $f.Params -}}
  {{- if $p.IsBasic -}}
    {{ $p.Name }} := param.{{firstToUpper $p.TypeName }}(p, "{{ $p.Name}}")
  {{- else -}}
    {{ $p.Name }} := &{{ $.Package }}.{{base $p.TypeName}}{}
    param.Object(p, {{ $p.Name }})
  {{- end }}
  {{ end }}

  v, err:= {{$.Package}}.{{$f.Name}}({{pluckStrings $f.Params "Name" | join ","}})
  Respond(w, r, v, err)
}
{{ end }}
{{ end }}

```

run `go generate order.go`, then http handler is ok.
```go
// generate by jia
package routes

import (
	"net/http"

	"somepkg/models"

	"github.com/relaxgo/tangram/param"
)

func FindOrderById(w http.ResponseWriter, r *http.Request) {
	p := NewRequestValue(r)

	userid := param.Int(p, "userid")
	orderid := param.Int(p, "orderid")

	v, err := models.FindOrderById(userid, orderid)
	Respond(w, r, v, err)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	p := NewRequestValue(r)

	userid := param.Int(p, "userid")
	title := param.String(p, "title")
	price := param.Float64(p, "price")

	v, err := models.CreateOrder(userid, title, price)
	Respond(w, r, v, err)
}
```

Of course, you can write some other template and generate the file you need.

## License
MIT
