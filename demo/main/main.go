package main

import (
	"fmt"
	"io/ioutil"

	"github.com/poldi1405/GoyEnc"
)

func main() {
	// asciiStr := "ABC"
	//asciiBytes := []byte(asciiStr)
	//reader := yenc.NewyEncReader([]byte("Hello world, I don't care for Datascience.\nABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz|\\/*-+&%$\"!?="), 128)
	file, err := ioutil.ReadFile("spec.txt")
	if err != nil {
		panic(err)
	}

	reader := yenc.NewyEncReader(file, 128)
	fmt.Println("File of length", len(file), "has been read")

	for {
		res, err := reader.ReadLine()
		fmt.Println(string(res))
		if err != nil {
			break
		}
	}

}
