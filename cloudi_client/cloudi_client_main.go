package main

import (
	"fmt"
	"os"
)

func main(){
	fileToFetch := os.Args[1]
	statusMsg := fetchFile(fileToFetch)
	if statusMsg != "" {
		fmt.Println(statusMsg)
	}
}