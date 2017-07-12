//go:generate jia -format=go -file $GOFILE -out ../routes/  -tpl ../route.tmpl
package models

type User struct {
	Id   int
	Name string
}
