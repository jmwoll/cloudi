package main

import (
	"fmt"
  "io"
  "io/ioutil"
	"net"
	"path/filepath"
	"os"
  "strconv"
  "strings"
  "crypto/sha1"
)

// Must match with client 
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
    // go sendFileToClient(connection)
    go requestHandler(connection)
	}
}

func requestHandler(connection net.Conn) {
  requestTypeBytes := make([]byte, 512)
  connection.Read(requestTypeBytes)
  requestType := strings.Trim(string(requestTypeBytes),":")
  fmt.Println("RequestType="+requestType)
  if requestType == "fetchFile" {
    sendFileToClient(connection)
	}
	if requestType == "listAllFiles" {
		sendAllFilesInformation(connection)
	}
}


func sendAllFilesInformation(connection net.Conn) {
  fmt.Println("Sending allFilesInformation")
	filesList := allFilesInCurrentDir()
	filesListStr := strings.Join(filesList, ":")
	fmt.Println(filesListStr)
	filesListBytes := []byte(filesListStr)
	filesListBytesLen := string(strconv.FormatInt(int64(len(filesListBytes)),10))//fillString(string(len(filesListBytes)),512)
	fmt.Println("~~~~~"+filesListBytesLen)
	// send info about byte size to client
	connection.Write([]byte(fillString(filesListBytesLen,512)))
	sendInChunks(filesListBytes, connection)
}


func sendInChunks(sourceBytes []byte, connection net.Conn){
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	inBufferIdx := 0
	for _,sourceByte := range sourceBytes {
			if inBufferIdx == BUFFERSIZE {
					connection.Write(sendBuffer)
					sendBuffer = make([]byte, BUFFERSIZE)
			} 
			inBufferIdx = inBufferIdx % BUFFERSIZE
			sendBuffer[inBufferIdx] = sourceByte
			inBufferIdx += 1
	}
	// Consider possibility of partially full buffer
	connection.Write(sendBuffer)
}



func allFilesInCurrentDir() []string {
	var filesList []string
	err := filepath.Walk(".",
    func(path string, info os.FileInfo, err error) error {
			if path == "." {
				// skip the reference to cwd
				// as it is annoying for fuzzy searches... 
				return nil
			}
			if err != nil {
					return err
			}
			filesList = append(filesList, path)
			//fmt.Println(path, info.Size())
			return nil
		})
		if err != nil {
				fmt.Println(err)
		}
		return filesList
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
