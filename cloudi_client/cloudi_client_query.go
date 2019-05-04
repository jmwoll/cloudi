package main

import (
    "fmt"
    "net"
    "strings"
)

func fuzzyStringMatch (){

}




// See http://rosettacode.org/wiki/Levenshtein_distances
func levenshteinDistance(s, t string) int {
    d := make([][]int, len(s)+1)
    for i := range d {
        d[i] = make([]int, len(t)+1)
    }
    for i := range d {
        d[i][0] = i
    }
    for j := range d[0] {
        d[0][j] = j
    }
    for j := 1; j <= len(t); j++ {
        for i := 1; i <= len(s); i++ {
            if s[i-1] == t[j-1] {
                d[i][j] = d[i-1][j-1]
            } else {
                min := d[i-1][j]
                if d[i][j-1] < min {
                    min = d[i][j-1]
                }
                if d[i-1][j-1] < min {
                    min = d[i-1][j-1]
                }
                d[i][j] = min + 1
            }
        }
 
    }
    return d[len(s)][len(t)]
}






func  listAllFiles (serverAddress string) ([]string,string) {
	//[]string rslt = nil
    // e.g. serverAddress := "localhost:27001"
    statusCode := "" /* Empty string signals success */
    connection, err := net.Dial("tcp", serverAddress)
    if err != nil {
        //panic(err)
        statusCode += "Error trying to connect to server"
        statusCode += serverAddress
        statusCode = err.Error()
        return nil,statusCode
    }
    defer connection.Close()
    // --- We need to inform the server that we are interested
    // --- in a file fetch action, because from perspective
    // --- of the server we could also perform another action.
    actionType := "listAllFiles"
    sendActionType(actionType,connection)
    // The server will now first send us the length of the
    // returned answer bytes, so we know what to expect.
    expectedAnswerLen := make([]byte,512)
    connection.Read(expectedAnswerLen)
    expectedAnswerLenStr := strings.Trim(string(expectedAnswerLen),":")
    fmt.Println("##############"+expectedAnswerLenStr)
	return nil,statusCode //rslt
}






