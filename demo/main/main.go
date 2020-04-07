package main

import (
	"fmt"
	"github.com/poldi1405/GoyEnc"
)

func main() {
	// asciiStr := "ABC"
	//asciiBytes := []byte(asciiStr)
	//reader := yenc.NewyEnc([]byte("Hello world, I don't care for Datascience.\nABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz|\\/*-+&%$\"!?="), 128)
	yEnc := yenc.NewyEnc()
	result := make(chan []byte)
	yEnc.EncodeFile("spec.txt", result)

	for encoded := range result {
		fmt.Print(encoded)
	}
}
