package main

import (
	"os"
)

func main(){
	fileToFetch := os.Args[1]
	fetchFile(fileToFetch)
}