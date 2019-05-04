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
	"net"
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
	return ""
}

