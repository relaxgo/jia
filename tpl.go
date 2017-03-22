package jia

import (
	"html/template"
	"reflect"
	"unicode"

	"github.com/Masterminds/sprig"
)

var StringsFuncs = template.FuncMap{
	"firstToUpper": func(str string) string {
		for i, v := range str {
			return string(unicode.ToUpper(v)) + str[i+1:]
		}
		return ""
	},
	"firstToLower": func(str string) string {
		for i, v := range str {
			return string(unicode.ToLower(v)) + str[i+1:]
		}
		return ""
	},
	"pluckStrings": func(src interface{}, fieldName string) []string {
		v := reflect.ValueOf(src)
		t := v.Type()
		if t.Kind() != reflect.Slice {
			panic("pluck need slice")
		}
		l := v.Len()
		s := make([]string, l, l)
		for i := 0; i < l; i++ {
			s[i] = v.Index(i).FieldByName(fieldName).String()
		}
		return s
	},
}

var BaseFuncs template.FuncMap

func init() {
	BaseFuncs = sprig.FuncMap()
	for k, v := range StringsFuncs {
		BaseFuncs[k] = v
	}
}
