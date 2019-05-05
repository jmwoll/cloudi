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
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

func pushFile(fileToPush,serverAddress string) string {
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
    // --- We need to inform the server that we are interested
    // --- in a file fetch action, because from perspective
    // --- of the server we could also perform another action.
    actionType := "pushFile"
	sendActionType(actionType,connection)
	// --- Now, we read in the contents of the file that we
	// --- want to add:
	fileBytes,err := ioutil.ReadFile(fileToPush)
	if err != nil {
		statusCode += "Aborting. Could not find file:\t"+fileToPush
		statusCode += err.Error()
		return statusCode
	}
	// --- The server does not know the size of the file, so
	// --- we need to inform the server about the amount of
	// --- bytes he has to receive.
	fileBytesLen := string(strconv.FormatInt(int64(len(fileBytes)),10))//fillString(string(len(filesListBytes)),512)
	fmt.Println("~~~~~"+fileBytesLen)
	// send info about byte size to client
	connection.Write([]byte(fillString(fileBytesLen,512)))
	// --- Then, we send the actual bytes of the file
	sendInChunks(fileBytes,connection)
	// --- Now, we want the server to give us an OK about
	// --- the upload: Did the file arrive at the server intact?
	// --- Thus, we send the server the sha1 hash of the file:
	
	return ""
}

