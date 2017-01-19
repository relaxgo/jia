package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
	file, err := os.Open(GoFilePath)
	defer file.Close()
	handleErr(err)
	data := Gen(file)
	if Output == "" {
		fmt.Println(string(data))
		return
	}
	err = ioutil.WriteFile(Output, data, 0640)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
