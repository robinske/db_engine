package main

import (
	"fmt"
	"log"
	//"io"
	"os"
	"bufio"
	//"strings"
)

func getInput() {
	in := bufio.NewReader(os.Stdin)
    line, err := in.ReadString('\n')// only accepts the first word with &message, accepts/concatenates words with just message)
	if err != nil {
		log.Fatal(err)
	}
    return line
}

func main() {
	msg := getInput()
	fmt.Printf("%s", msg)
	
}