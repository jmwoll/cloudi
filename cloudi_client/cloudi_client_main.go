package main

import (
	"fmt"
	"os"
)

func main(){
	serverAddress := "localhost:27001"
	fileToFetch := os.Args[1]
	statusMsg := fetchFile(fileToFetch,serverAddress)
	if statusMsg != "" {
		fmt.Println(statusMsg)
	}
}