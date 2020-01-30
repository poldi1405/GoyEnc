package main

import (
	"fmt"
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

	reader := yenc.NewyEnc(os.Args[1], 128, true)

	reader.EncodeFile()
}
