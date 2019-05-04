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

const BUFFERSIZE = 1024

func main() {
	server, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go sendFileToClient(connection)
	}
}


func sendFileToClient(connection net.Conn) {
	fmt.Println("A client has connected!")
  defer connection.Close()
  // --- We want to read in the file name
  // --- the client wants to receive from
  // --- the connection.
  bufferFileName := make([]byte, 64)
  connection.Read(bufferFileName)
  fmt.Println("Client Query:")
  fileNameFromClient := strings.Trim(string(bufferFileName), ":")
  fmt.Println(fileNameFromClient)
  // ---
	file, err := os.Open(fileNameFromClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 512)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
  fmt.Println("File has been sent!")
  h := sha1.New()
  fileBytes,err := ioutil.ReadFile(fileNameFromClient)
  h.Write(fileBytes)
  hashSum := h.Sum(nil)
  fmt.Println("Server side hashsum")
  fmt.Printf("%x",hashSum)
  connection.Write(hashSum)
	return
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
