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
	"os"
)

func main() {
	serverAddress := "localhost:27001"
	actionType := os.Args[1]
	if actionType == "get" {
		fileToFetch := os.Args[2]
		statusMsg := fetchFile(fileToFetch, serverAddress)
		if statusMsg != "" {
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)
		}
	}
	if actionType == "list" {
		allFiles, statusMsg := listAllFiles(serverAddress)
		if statusMsg != "" {
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)
		}
		fmt.Print(ansiColor("blue"))
		for _, file := range allFiles {
			fmt.Println(file)
		}
		//fmt.Println(allFiles)
	}
	if actionType == "find" {
		fileQuery := os.Args[2]
		ratio, fileFound, statusMsg := findFile(fileQuery, serverAddress)
		if statusMsg != "" {
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)
		}
		fmt.Printf("levensthein=%f", ratio)
		if ratio > levenstheinMinimumRatio() {
			fmt.Print(ansiColor("green"))
			fmt.Println("found file and copied it into cwd:\t" + fileFound)
		} else {
			fmt.Print(ansiColor("yellow"))
			fmt.Println("could not find file:\t" + fileQuery)
		}
	}
	if actionType == "add" {
		fileToAdd := os.Args[2]
		statusMsg := pushFile(fileToAdd, serverAddress)
		if statusMsg != "" {
			fmt.Print(ansiColor("red"))
			fmt.Println(statusMsg)
		}
	}

}
