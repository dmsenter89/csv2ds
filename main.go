package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	files := os.Args[1:]
	for _, f := range files {
		fmt.Printf("%s\n", filenameWithoutExtension(f))
	}

	fmt.Println("Initial commit line.")
}

func usage() {
	fmt.Printf("Usage: %s file1 [file2...]\n", path.Base(os.Args[0]))
}

func filenameWithoutExtension(filepath string) string {
	var fileName string = path.Base(filepath)
	var fileExtension string = path.Ext(filepath)
	return strings.TrimSuffix(fileName, fileExtension)
}
