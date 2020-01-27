package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/poldi1405/GoyEnc"
)

func main() {
	// asciiStr := "ABC"
	//asciiBytes := []byte(asciiStr)
	//reader := yenc.NewyEnc([]byte("ABCÖÜßöäü"), 128)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please enter a Filename")
		os.Exit(2)
	}

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	reader := yenc.NewyEnc(file, 128, true)
	fmt.Fprintln(os.Stderr, "File of length", len(file), "has been read")

	for {
		res, err := reader.ReadLine()
		fmt.Print(string(res) + "\r\n")
		if err != nil {
			break
		}
	}

}
