package main

import (
	"io"
	"log"
	"os"
	P "spalk/test/parser"
)

func main() {
	/*
	Purpose of this portion of the code is to get the input from the file
	*/
	data, err := io.ReadAll(os.Stdin)

	var lines string
	if err != nil {
		log.Fatal(err)
	}
	lines = string(data)

	dataStream := []byte(lines)

	P.ParseMPEG(dataStream)
}