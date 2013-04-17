package main

import (
	"os"
	"io"
	"db_engine/client"
	"log"
	"bufio"
	"strings"
)


func readFile(file string) string {
    fo, err := os.OpenFile(file, os.O_RDWR, 0666)
    if err != nil {
        log.Fatal(err)
    }
    
    defer fo.Close()

	r := bufio.NewReader(fo)

	lineList := []string{}

	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		s := string(l)

		lineList = append(lineList, s)
	}

	s := strings.Join(lineList, ";")

	return s
}


func main() {
	
	fileText := readFile("working/students.txt")
	client.Connect(fileText)

}