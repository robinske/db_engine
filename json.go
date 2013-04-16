package main

import (
    // "net"
    "log"
    // "io"
    // "os"
    // "strings"
    "fmt"
    "encoding/json"
)

type Dictionary map[string]interface{}

type Object struct {
    Dict Dictionary
}


func main() {

	key := "thekey"
	value := Dictionary{"value as key":true}
	var o = Object {}
	o.Dict = Dictionary{key:value}

	b, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("normal dictionary: %v\n", o)

	fmt.Printf("json BYTE encoded dictionary: %v\n", b)
	fmt.Printf("json STRING encoded dictionary: %v\n", string(b))

	var d = map[string]interface{}{}

	error := json.Unmarshal(b, &d)
	if error != nil {
		log.Fatal(err)
	}

	fmt.Printf("json DECODED dictionary: %v\n", d)

}