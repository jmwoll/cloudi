package main

import (
	"fmt"
  "io"
  "io/ioutil"
	"net"
	"os"
	"strconv"
  "strings"
  "crypto/sha1"
)

// Have a look at https://mrwaggel.be/post/golang-transfer-a-file-over-a-tcp-socket/
// Needs to be same as in server...
const BUFFERSIZE = 1024

func fetchFile(fileNameQueryArg string) {
	connection, err := net.Dial("tcp", "localhost:27001")
	if err != nil {
		panic(err)
	}
  defer connection.Close()
  // --- Modified the server to receive the filename from 
  // --- the client. So the client now needs to send
  // --- the filename of interest to server. Let's do this here:
  fileNameQuery := fillString(fileNameQueryArg,512)
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
	bufferFileName := make([]byte, 512)
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
  shaHashServer := make([]byte, 20)
  connection.Read(shaHashServer)
  fmt.Println("The sha1 hash of the received file:")
  fmt.Printf("%x",shaHashServer)
  // --- Now we can compar the sha1 hash that
  // --- the file should have based on the info 
  // --- above with the sha1 hash we actually
  // --- compute for the fetched file on the client
  h := sha1.New()
  fileBytes,err := ioutil.ReadFile(fileName)
  h.Write(fileBytes)
  hashSumClient := h.Sum(nil)
  fmt.Println("actual client hash:")
  fmt.Printf("%x",hashSumClient)
  hashSumClientAsByteArray := make([]byte,20)
  for idx,b := range hashSumClient {
    hashSumClientAsByteArray[idx] = b
  }
  if byteArraysEqual(shaHashServer,hashSumClientAsByteArray) {
    fmt.Println("sha1 hash matches")
  }else{
    fmt.Println("sha1 hash do not match:")
    fmt.Printf("%x (on client) != %x (on server)", hashSumClientAsByteArray, shaHashServer)
  }
}


func byteArraysEqual(as,bs []byte) bool {
  for idx,a := range as {
    if a != bs[idx] {
      return false
    }
  }
  return true
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