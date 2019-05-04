package main

import (
	"fmt"
	"os"
)

func main(){
	serverAddress := "localhost:27001"
	actionType := os.Args[1]
	if actionType == "get" {
		fileToFetch := os.Args[2]
		statusMsg := fetchFile(fileToFetch,serverAddress)
		if statusMsg != "" {
			fmt.Println(statusMsg)
		}
	}
	if actionType == "list" {
		allFiles,statusMsg := listAllFiles(serverAddress)
		if statusMsg != "" {
			fmt.Println(statusMsg)
			fmt.Println(allFiles)
		}
	} 
}