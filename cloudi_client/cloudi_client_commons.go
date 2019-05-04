package main

import (
	"fmt"
)

func showDebugMessages () bool {
	return false
}

func debug(msg string){
	if showDebugMessages() {
		fmt.Println(msg)
	}
}