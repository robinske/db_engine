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

type StringDict map[string]string
type FlexDict map[string]interface{}

type Object struct {
    Dict FlexDict
}

func main() {

	key := "thekey"
	key2 := "value as key"
	value2 := []int{2,3,4,5}
	value := FlexDict{key2:value2}
	// var o = Object {} // instantiate the object
	// o.Dict = FlexDict{key:value} // give the method Dict a value

	// b, err := json.Marshal(o)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("normal dictionary: %v\n", o)

	// fmt.Printf("json BYTE encoded dictionary: %v\n", b)
	// fmt.Printf("json STRING encoded dictionary: %v\n", string(b))

	// var d = FlexDict {}

	// error := json.Unmarshal(b, &d)
	// if error != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("json DECODED dictionary: %v\n", d)

/// NOW WITHOUT OBJECTS
	t := 5
	var t = 5

	var t = FlexDict {}
	t[key] = value

	j, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NOW WITH JUST TYPES")

	fmt.Printf("normal dictionary: %v\n", t)

	fmt.Printf("json BYTE encoded dictionary: %v\n", j)
	fmt.Printf("json STRING encoded dictionary: %v\n", string(j))

	// get at the values
	// fmt.Printf("value: %v", j[key])

	var u = FlexDict {}

	newerror := json.Unmarshal(j, &u)
	if newerror != nil {
		log.Fatal(err)
	}

	fmt.Printf("json DECODED dictionary: %v\n", u)



}