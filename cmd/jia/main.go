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
)

func main() {
	flag.Parse()
	if *GoFilePath == "" {
		flag.Usage()
		return
	}

	file, err := os.Open(*GoFilePath)
	defer file.Close()
	handleErr("failed to open file ", err)

	f, err := jia.Parse(path.Base(*GoFilePath), file)
	handleErr("fialed to parse go file", err)
	data, err := Render(f, tpls.EchoTemp)
	handleErr("failed to generate go file", err)

	if *Output != "" {
		err = ioutil.WriteFile(*Output, data, 0640)
		handleErr("faild to write file", err)
	} else {
		fmt.Println(string(data))
	}
}

func handleErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}

func Render(f *jia.GoFile, t *template.Template) ([]byte, error) {
	bf := bytes.NewBuffer(nil)
	err := t.Execute(bf, f)
	if err != nil {
		return nil, err
	}
	return format.Source(bf.Bytes())
}
