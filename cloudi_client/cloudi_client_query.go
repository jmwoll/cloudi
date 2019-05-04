package main

import (
    "fmt"
    "math"
    "net"
    "strings"
    "strconv"
)

func fuzzyStringMatch (){

}


// see http://www.golangprograms.com/golang-program-for-implementation-of-levenshtein-distance.html
/*func levenshteinDistance2(argStr1, argStr2 string) int {
    str1 := []rune(argStr1)
    str2 := []rune(argStr2)
    s1len := len(str1)
    s2len := len(str2)
    column := make([]int, len(str1)+1)
 
    for y := 1; y <= s1len; y++ {
        column[y] = y
    }
    for x := 1; x <= s2len; x++ {
        column[0] = x
        lastkey := x - 1
        for y := 1; y <= s1len; y++ {
            oldkey := column[y]
            var incr int
            if str1[y-1] != str2[x-1] {
                incr = 1
            }
 
            column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
            lastkey = oldkey
        }
    }
    return column[s1len]
}
 
func minimum(a, b, c int) int {
    if a < b {
        if a < c {
            return a
        }
    } else {
        if b < c {
            return b
        }
    }
    return c
}
*/



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

func levenstheinMinimumRatio() float64 {
    return 0.25
}

func findFile (toFind, serverAddress string) (float64,string,string) {
    ratio := 0.0
    allFiles,statusCode := listAllFiles(serverAddress)
    if statusCode != "" {
        return 0.0,"",statusCode
    }
    matchedFile := ""
    minDist := 5000 /* assume str len << 5000 */
    for _,file := range(allFiles){
        // TODO: file will probably be a full path,
        // and as such consist of directory+file.
        // Then we probably have to split at path sep
        // to only compare the file name without dir name.     
        curDist := levenshteinDistance(file,toFind)
        maxL := math.Max(float64(len(file)),float64(len(toFind)))
        ratio = float64(curDist) / maxL
        ratio = 1.0 - ratio
        if curDist < minDist && ratio > levenstheinMinimumRatio() {
            minDist = curDist
            matchedFile = file
        }
    }
    fmt.Println("toFind:"+toFind)
    fmt.Println("THE MATCHED FILE IS:"+matchedFile)
    if matchedFile != "" {
        statusCode += fetchFile(matchedFile, serverAddress)
    }
    return ratio,matchedFile,statusCode
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
    expectedAnswerLenStr := strings.Trim(strings.Trim(string(expectedAnswerLen),":"),"\x00")
    fmt.Println("##############"+expectedAnswerLenStr)
    numBytes, err := strconv.Atoi(expectedAnswerLenStr)
    if err != nil {
        statusCode += err.Error()
    }
    listOfFiles := make([]byte,numBytes)
    connection.Read(listOfFiles)
    listOfFilesStr := string(listOfFiles)
    //fmt.Println("|||| =>"+listOfFilesStr)
    rslt := strings.Split(listOfFilesStr,":")
	return rslt,statusCode //rslt
}






