package main

import (
	"fmt"
	"strings"
)

func parseRequest(msg string) {
	msgSplit := strings.Split(msg, " ")

	fmt.Printf("%q\n", msg)
	fmt.Printf("%q\n", msgSplit)

	if msgSplit[0] == "GET" {
		// GET function
		fmt.Printf("you called the %s function\n", msgSplit[0])
	}

	if msgSplit[0] == "PUT" {
		// PUT function
		fmt.Printf("you called the %s function\n", msgSplit[0])
	}

	if msgSplit[0] == "SAVE" {
		// WRITE TO FILE function
		fmt.Printf("you called the %s function\n", msgSplit[0])
	}
	
    // fmt.Printf("%q\n", &message)
	
}

func main() {
	msg := "GET some other things"
	parseRequest(msg)	
}

	// fmt.Printf("%q\n", list[1])
	// fmt.Printf("%q\n", strings.Split(abc, ","))
	// fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
	// fmt.Printf("%q\n", strings.Split(" xyz ", ""))
	// fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))