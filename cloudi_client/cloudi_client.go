package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Have a look at https://mrwaggel.be/post/golang-transfer-a-file-over-a-tcp-socket/
// Needs to be same as in server...
const BUFFERSIZE = 1024

func main() {
	connection, err := net.Dial("tcp", "localhost:27001")
	if err != nil {
		panic(err)
	}
  defer connection.Close()
  // --- Modified the server to receive the filename from 
  // --- the client. So the client now needs to send
  // --- the filename of interest to server. Let's do this here:
  fileNameQuery := fillString(os.Args[1],64)
  fmt.Println("fileNameQuery="+fileNameQuery)
  connection.Write([]byte(fileNameQuery))
  // ---
  // ---
  // --- Although we send the fileName to the server already, the
  // --- will still likely return *another file* to the client!!
  // --- Because nobody can remember the filename 100% correctly,
  // --- the server should perform fuzzy matching on the filename
  // --- request coming from the client. This way, the user wont
  // --- stumble over his own feet if he has a single typo...
  // ---
	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	
	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	
	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	
	newFile, err := os.Create(fileName)
	
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64
	
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}


func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}