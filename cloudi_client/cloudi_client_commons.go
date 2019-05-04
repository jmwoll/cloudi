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