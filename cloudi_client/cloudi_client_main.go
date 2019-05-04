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
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)
		}
	}
	if actionType == "list" {
		allFiles,statusMsg := listAllFiles(serverAddress)
		if statusMsg != "" {
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)	
		}
		fmt.Print(ansiColor("blue"))
		for _,file := range allFiles {
			fmt.Println(file)
		} 
		//fmt.Println(allFiles)
	} 
	if actionType == "find" {
		fileQuery := os.Args[2]
		statusMsg := findFile(fileQuery,serverAddress)
		if statusMsg != "" {
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)
		}
	}
}