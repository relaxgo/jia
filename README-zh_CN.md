## Jia
通过`go/ast`等库解析go文件的结构，然后用自定义的模版，生成所需要的其他文件，如API的controler, 前端的js文件

## Warning
目前虽然已经可以使用，但是还在完善中，接口随时可能会更改

## Example
直接把web API写成函数, 而不需要管参数怎么获得
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
添加你的路由函数的模版
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

运行`go generate order.go`, 然后就你的http handler 函数就有啦
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

当然也可以通过模版生成其他文件，比如前端的js文件，swagger文件

## License
MIT
