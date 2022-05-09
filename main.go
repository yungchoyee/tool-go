package main

import (
	"fmt"
	"os"
	"tool-go/file"
)

func main() {
	curPath, _ := os.Getwd()
	filenames, _ := file.GetFileNames(curPath)
	fmt.Println(filenames)
}
