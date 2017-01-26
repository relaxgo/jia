package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var (
	Output     string
	GoFilePath string
)

func init() {
	flag.StringVar(&Output, "o", "", "输出文件路径")
	flag.StringVar(&GoFilePath, "f", "", "go 文件")
}

func main() {
	flag.Parse()
	if GoFilePath == "" {
		flag.Usage()
		return
	}
	file, err := os.Open(GoFilePath)
	defer file.Close()
	handleErr(err)
	f, err := Parse(path.Base(GoFilePath), file)
	handleErr(err)
	data, err := Render(f, EchoTemp)
	handleErr(err)
	if Output != "" {
		err = ioutil.WriteFile(Output, data, 0640)
		handleErr(err)
	} else {
		fmt.Println(string(data))
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
