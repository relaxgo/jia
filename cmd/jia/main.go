package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/relaxgo/jia"
	"github.com/relaxgo/jia/tpls"
)

var (
	Output     = flag.String("o", "", "输出文件路径")
	GoFilePath = flag.String("f", "", "go 文件")
	Template   = flag.String("t", "", "模版文件")
	Format     = flag.Bool("format", false, "格式化")
)

func main() {
	flag.Parse()
	if *GoFilePath == "" {
		flag.Usage()
		return
	}

	file, err := os.Open(*GoFilePath)
	defer file.Close()
	handleErr("open file ", err)

	f, err := jia.Parse(*GoFilePath, file)
	handleErr("parse go file", err)

	t := LoadTemplate(*Template)
	data, err := Render(f, t)
	fmt.Println("success to  generate file")

	if *Output != "" {
		filePath := Resolve(*GoFilePath, *Output)
		err = ioutil.WriteFile(filePath, data, 0640)
		handleErr("write file", err)
	} else {
		fmt.Println(string(data))
	}
}

func LoadTemplate(f string) *template.Template {
	if f == "" {
		return tpls.EchoTemp
	}
	data, err := ioutil.ReadFile(f)
	handleErr("load template", err)
	t, err := template.New("root").Funcs(jia.BaseFuncs).Parse(string(data))
	handleErr("parse template", err)
	return t
}

func Render(f *jia.GoFile, t *template.Template) ([]byte, error) {
	bf := bytes.NewBuffer(nil)
	err := t.Execute(bf, f)
	if err != nil {
		handleErr("execute template", err)
		return nil, err
	}

	if *Format {
		v, err := format.Source(bf.Bytes())
		handleErr("format go file", err)
		return v, err
	}
	return bf.Bytes(), nil
}

func IsDir(p string) bool {
	info, err := os.Stat(p)
	// FIXME p may be not exist
	if err != nil {
		return false
	}
	return info.IsDir()
}

func Resolve(inPath, outPath string) string {
	if IsDir(outPath) {
		return path.Join(outPath, path.Base(inPath))
	}
	return outPath
}

func handleErr(task string, err error) {
	if err != nil {
		fmt.Println("failed to ", task)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("success to ", task)
	}
}
