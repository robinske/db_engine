package main

import (
	"fmt"
	"os"
	"bufio"
	"log"
	"io"
	"strings"
	"db_engine/client"
)

type Details struct {
	Name string
	Email string
	Twitter string
	Phone string
	Github string
}

type Student struct {
	Id int
	Details Details
}

type Students map[string]interface{}

func readFile(file string) {
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

	splitLines(lineList)
}

func splitLines(list []string) {
		

		for i := 0; i < len(list); i++ {

			elems := strings.Split(list[i], "\t")
			
			fmt.Println(elems)

			s := Student{elems[0], elems[1], elems[2], elems[3], elems[4]}

			// d := {id, s}		

		}

	

	b, err := json.Marshal(d)
	//fmt.Println(list)
}

func main() {
	readFile("students.txt")

	// needs to return JSON

}