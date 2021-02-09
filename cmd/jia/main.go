package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/relaxgo/jia"
)

var (
	output      = flag.String("out", "", "输出文件路径")
	goFilePath  = flag.String("file", "", "go 文件")
	filTemplate = flag.String("tpl", "", "模版文件")
	fileFormat  = flag.String("format", "", "格式化")

	infoLog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
)

const (
	formatGo = "go"
)

func main() {
	flag.Parse()
	if *goFilePath == "" || *filTemplate == "" {
		flag.Usage()
		return
	}

	file, err := os.Open(*goFilePath)
	defer file.Close()
	handleErr("open file ", err)

	abspath, err := filepath.Abs(*goFilePath)
	handleErr("expend go file path", err)
	f, err := jia.ParsePackage(abspath)
	handleErr("parse go file", err)

	t := loadTemplate(*filTemplate)
	data, err := render(f, t)
	handleErr("generate file", err)

	if *output != "" {
		filePath := resolve(*goFilePath, *output)
		err = ioutil.WriteFile(filePath, data, 0640)
		handleErr("write file", err)
	} else {
		fmt.Println(string(data))
	}
}

func loadTemplate(f string) *template.Template {
	if f == "" {
		panic("need template file")
	}
	data, err := ioutil.ReadFile(f)
	handleErr("load template", err)
	t, err := template.New("root").Funcs(jia.BaseFuncs).Parse(string(data))
	handleErr("parse template", err)
	return t
}

func render(f *jia.GoFile, t *template.Template) ([]byte, error) {
	bf := bytes.NewBuffer(nil)
	err := t.Execute(bf, f)
	if err != nil {
		handleErr("execute template", err)
		return nil, err
	}

	if *fileFormat == formatGo {
		v, err := format.Source(bf.Bytes())
		handleErr("format go file", err)
		return v, err
	}
	return bf.Bytes(), nil
}

func isDir(p string) bool {
	info, err := os.Stat(p)
	// FIXME p may be not exist
	if err != nil {
		return false
	}
	return info.IsDir()
}

func resolve(inPath, outPath string) string {
	if isDir(outPath) {
		return path.Join(outPath, path.Base(inPath))
	}
	return outPath
}

func handleErr(task string, err error) {
	if err != nil {
		infoLog.Println("failed to ", task)
		infoLog.Println(err)
		os.Exit(1)
	} else {
		infoLog.Println("success to ", task)
	}
}
