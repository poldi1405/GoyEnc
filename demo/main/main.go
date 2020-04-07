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

	yEnc := yenc.NewyEnc()
	result := make(chan []byte)
	yEnc.EncodeFile(os.Args[1], result)

	for encoded := range result {
		fmt.Print(encoded)
	}
}
