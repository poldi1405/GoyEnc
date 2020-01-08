package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/poldi1405/GoyEncode"
)

func main() {
	// asciiStr := "ABC"
	//asciiBytes := []byte(asciiStr)
	//reader := yenc.NewyEncReader("Hello world, I don't care for Datascience.\nABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz|\\/*-+&%$\"!?=", 128)
	file, err := ioutil.ReadFile("spec.txt")
	if err != nil {
		panic(err)
	}
	reader := yenc.NewyEncReader(string(file), 128)
	p := make([]byte, 4*2)
	for {
		n, linebreak, err := reader.Read(p)
		if linebreak {
			fmt.Print("DINGDING\n")
		}
		if err == io.EOF {
			break
		}
		if n >= len(p) {
			n--
		}
		fmt.Print(string(p[:n]))
	}
}
