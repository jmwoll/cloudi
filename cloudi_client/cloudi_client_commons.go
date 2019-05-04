package main

import (
	"fmt"
)

func showDebugMessages () bool {
	return false
}

func debugMsg(msg string){
	if showDebugMessages() {
		fmt.Println(msg)
	}
}