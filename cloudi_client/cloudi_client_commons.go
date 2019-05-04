package main

import (
	"fmt"
	"net"
)

func showDebugMessages () bool {
	return false
}

func debugMsg(msg string){
	if showDebugMessages() {
		fmt.Println(msg)
	}
}


func sendActionType(actionType string, connection net.Conn){
  actionType = fillString(actionType,512)
  actionTypeBytes := []byte(actionType)
  actionTypeBytesPadded := make([]byte,512)
  for index,b := range actionTypeBytes {
    actionTypeBytesPadded[index] = b
  }
  debugMsg("sending::::"+string(actionTypeBytesPadded))
  connection.Write(actionTypeBytesPadded)
}

/*
Black        0;30     Dark Gray     1;30
Red          0;31     Light Red     1;31
Green        0;32     Light Green   1;32
Brown/Orange 0;33     Yellow        1;33
Blue         0;34     Light Blue    1;34
Purple       0;35     Light Purple  1;35
Cyan         0;36     Light Cyan    1;36
Light Gray   0;37     White         1;37
No Color 0
*/
func ansiColor(colorName string) string {
	aCols := map[string]string {
		"black": "0;30", "red": "0;31",
		"blue": "0;34", "yellow": "1;33",
		"green": "0;32", "purple": "0;35",
	}
	if val, ok := aCols[colorName]; ok {
		return "\033["+val+"m"
	}
	return ""
}

