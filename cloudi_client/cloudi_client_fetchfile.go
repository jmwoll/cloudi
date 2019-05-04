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


// FetchFile mercilessly tries to fetch a file
// with the supplied clearname <fileNamyQueryArg>
// from the server. If the file with *exactly this name*
// does not exist on the server, an error message
// will be returned as a string.
// If the file is found on the server, it is copied
// into the current working directory on the client.
// Then, a sha1 hash checksum check is performed to
// assert that the file has been correctly transferred
// from the server to the client. If the check sums
// don't match, a meaningful error message is returned
// as a string. 
// This method guarantees that if the error message string
// returned *is empty*, then the client now has a
// correct copy of the target file in the current working dir.
func fetchFile(fileNameQueryArg,serverAddress string) string {
  // e.g. serverAddress := "localhost:27001"
  statusCode := "" /* Empty string signals success */
	connection, err := net.Dial("tcp", serverAddress)
	if err != nil {
    //panic(err)
    statusCode += "Error trying to connect to server"
    statusCode += serverAddress
    statusCode = err.Error()
    return statusCode
  }
  defer connection.Close()
  // --- Modified the server to receive the filename from 
  // --- the client. So the client now needs to send
  // --- the filename of interest to server. Let's do this here:
  fileNameQuery := fillString(fileNameQueryArg,512)
  debugMsg("fileNameQuery="+fileNameQuery)
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
	debugMsg("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 512)
	bufferFileSize := make([]byte, 10)
	
	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	
	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	
	newFile, err := os.Create(fileName)
	
	if err != nil {
    //panic(err)
    statusCode += "Error: file '"+strings.Trim(fileNameQuery,":")+"' not found on server\n"
    statusCode += fileName
    statusCode += err.Error()
    return statusCode
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
  debugMsg("Received file completely!")
  shaHashServer := make([]byte, 20)
  connection.Read(shaHashServer)
  debugMsg("The sha1 hash of the received file:")
  debugMsg(fmt.Sprintf("%x",shaHashServer))
  // --- Now we can compar the sha1 hash that
  // --- the file should have based on the info 
  // --- above with the sha1 hash we actually
  // --- compute for the fetched file on the client
  h := sha1.New()
  fileBytes,err := ioutil.ReadFile(fileName)
  h.Write(fileBytes)
  hashSumClient := h.Sum(nil)
  debugMsg("actual client hash:")
  debugMsg(fmt.Sprintf("%x",hashSumClient))
  hashSumClientAsByteArray := make([]byte,20)
  for idx,b := range hashSumClient {
    hashSumClientAsByteArray[idx] = b
  }
  if byteArraysEqual(shaHashServer,hashSumClientAsByteArray) {
    debugMsg("sha1 hash matches")
  }else{
    debugMsg("sha1 hash do not match:")
    debugMsg(fmt.Sprintf("%x (on client) != %x (on server)",
     hashSumClientAsByteArray, shaHashServer))
    statusCode += "Error receiving file from server:\n"
    statusCode += "sha1 hashes do not match:"
    statusCode += fmt.Sprintf("%x (on client) != %x (on server)", hashSumClientAsByteArray, shaHashServer)
  }
  return statusCode
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