package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	files := os.Args[1:]
	for _, f := range files {
		fmt.Printf("%s\n", filenameWithoutExtension(f))
	}

	fmt.Println("Initial commit line.")
}

func filenameWithoutExtension(filepath string) string {
	var fileName string = path.Base(filepath)
	var fileExtension string = path.Ext(filepath)
	return strings.TrimSuffix(fileName, fileExtension)
}
