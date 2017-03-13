package jia

import (
	"html/template"
	"path"
	"reflect"
	"strings"
	"unicode"
)

var BaseFuncs = template.FuncMap{
	"isBaseType": func(str string) bool {
		return true
	},
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
	"toLower": func(s string) string {
		return strings.ToLower(s)
	},
	"base": func(s string) string {
		return path.Base(s)
	},
	"dir": func(s string) string {
		return path.Dir(s)
	},
	"joinField": func(slice interface{}, fieldName, sep string) string {
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
}
