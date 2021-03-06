/*
 * This file is part of the cloudi project (https://github.com/jmwoll/cloudi/).
 * Copyright (c) 2019 Jan Wollschläger.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	requestType := strings.Trim(string(requestTypeBytes), ":")
	fmt.Println("RequestType=" + requestType)
	if requestType == "fetchFile" {
		sendFileToClient(connection)
	}
	if requestType == "listAllFiles" {
		sendAllFilesInformation(connection)
	}
	if requestType == "pushFile" {
		getFileFromClient(connection)
	}
}

func getFileFromClient(connection net.Conn) {
	fmt.Println("getting file from client")
	// --- First, we want to retrieve the
	// --- name of the file from the client:
	nameOfFileToAddBytes := make([]byte, 512)
	connection.Read(nameOfFileToAddBytes)
	nameOfFile := strings.Trim(string(nameOfFileToAddBytes), ":")
	fmt.Println("So you want to add the file:\t" + nameOfFile)
	expectedFileSizeBytes := make([]byte, 512)
	connection.Read(expectedFileSizeBytes)
	//expectedFileSizeStr := strings.Trim(string(expectedFileSize), ":")
	fileSize, _ := strconv.ParseInt(strings.Trim(string(expectedFileSizeBytes), ":"), 10, 64)
	fileBytes := make([]byte, fileSize)
	// --- First step: get sha1 hash from client
	shaOnClient := make([]byte, 20)
	connection.Read(shaOnClient)
	fmt.Print("The sha1 hash of the received file on client:")
	fmt.Printf("%x", shaOnClient)
	// --- Now that we know the number of bytes
	// --- to receive, let's get the file's bytes:
	connection.Read(fileBytes)
	// --- In fileBytes we now have the contents
	// --- of the files we wish to add to the server.
	// --- Now, we need to check that the sha1 sum matches,
	// --- and only then store the file on the server.
	// --- Second step: compute sha1 hash on server
	h := sha1.New()
	h.Write(fileBytes)
	hashSum := h.Sum(nil)
	fmt.Println("Server side hashsum")
	fmt.Printf("%x", hashSum)
	connection.Write(hashSum)
	successfulPush := byteArraysEqual(hashSum, shaOnClient)
	if successfulPush {
		fmt.Println("Push on server successful")
		err := ioutil.WriteFile(nameOfFile, fileBytes, 0644)
		if err != nil {
			fmt.Println("Error writing file:")
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Push on server failed. Sha1 hashes do not match:")
		fmt.Printf("%x (on client) != %x (on server)", shaOnClient, hashSum)
	}
}

func sendAllFilesInformation(connection net.Conn) {
	fmt.Println("Sending allFilesInformation")
	filesList := allFilesInCurrentDir()
	filesListStr := strings.Join(filesList, ":")
	fmt.Println(filesListStr)
	filesListBytes := []byte(filesListStr)
	filesListBytesLen := string(strconv.FormatInt(int64(len(filesListBytes)), 10)) //fillString(string(len(filesListBytes)),512)
	fmt.Println("~~~~~" + filesListBytesLen)
	// send info about byte size to client
	connection.Write([]byte(fillString(filesListBytesLen, 512)))
	sendInChunks(filesListBytes, connection)
}

func sendInChunks(sourceBytes []byte, connection net.Conn) {
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	inBufferIdx := 0
	for _, sourceByte := range sourceBytes {
		if inBufferIdx == BUFFERSIZE {
			connection.Write(sendBuffer)
			sendBuffer = make([]byte, BUFFERSIZE)
		}
		inBufferIdx = inBufferIdx % BUFFERSIZE
		sendBuffer[inBufferIdx] = sourceByte
		inBufferIdx += 1
	}
	// Consider possibility of partially full buffer
	// TODO: fix this part here. If we truncate the
	// send buffer to len(sourceBytes) % BUFFERSIZE
	// TODO: apply same fix to client code...
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
	fileBytes, err := ioutil.ReadFile(fileNameFromClient)
	h.Write(fileBytes)
	hashSum := h.Sum(nil)
	fmt.Println("Server side hashsum")
	fmt.Printf("%x", hashSum)
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

func byteArraysEqual(as, bs []byte) bool {
	for idx, a := range as {
		if a != bs[idx] {
			return false
		}
	}
	return true
}
