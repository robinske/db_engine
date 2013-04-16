package main

import (
    // "net"
    // "log"
    // "io"
    // "os"
    // "strings"
    "fmt"
    //"encoding/json"
)

type Dictionary map[string]string

type Object struct {
    Dict Dictionary
}

type Dict map[string]string
// get at something within the map
// func (o *Object) Value() string {
//     	return o.Dictionary[key]
// 	}

func main() {
	
	var d = Dict {}

	d["this"] = "that"
	d["another"] = "thing"

	fmt.Printf("print dict: %v\n", d)

	key := "thekey"
	value := "thevalue"

	var o = Object {}
	o.Dict = map[string]string{key:value}

	fmt.Printf("print OBJECT dict: %v\n", o)
	fmt.Printf("print OBJECT dict VALUE: %v\n", o.Dict[key])

}