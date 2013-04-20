package main

import (
	"fmt"
	//"db_engine/flat"
	"io/ioutil"
	"log"
	//"strings"
	"encoding/json"
	"strconv"
)

// type Dictionary map[string]interface{}
// type JSON map[string]interface{}

var lkey = ""

// var flattenedJSON = Dictionary {}

func flatten(inputJSON map[string]interface{}, lkey string, flattened *map[string]string) {
	for rkey, value := range inputJSON {
		key := lkey+rkey
		if _, ok := value.(string); ok {
			(*flattened)[key] = value.(string)
		} else if _, ok := value.([]interface{}); ok {
			// value is a list
			// flatten the values in the list
			for i := 0; i<len(value.([]interface{})); i++ {
				if _, ok := value.([]string); ok {
					fmt.Println("work?")
					stringI := string(i)
					(*flattened)[stringI] = value.(string)
					/// think this is wrong

				} else {
				flatten(value.([]interface{})[i].(map[string]interface{}), key+":"+strconv.Itoa(i)+":", flattened)
				}
			}
		} else {
			flatten(value.(map[string]interface{}), key+":", flattened)
			// 0 won't print anything but does contain a byte
			// maybe for testing use a unicode character
		}
	}
}

func defineType(value interface{}, key string, flattened *map[string]string) interface{} {
	if _, ok := value.(string); ok {
		(*flattened)[key] = value.(string)
	} else if _, ok := value.(int); ok {
		(*flattened)[key] = value.(int)
	} else if _, ok := value.(float64); ok {
		(*flattened)[key] = value.(float64)
	} else if _, ok := value.(bool); ok {
		(*flattened)[key] = value.(bool)
	} else {
		fmt.Println("RAWR")
	}
}





func load(filename string) {
    fileContents, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    mappedJSON := decodeJSON(fileContents)
    var flattened = make(map[string]string)
    flatten(mappedJSON, lkey, &flattened)
    for key, value := range flattened {
    	fmt.Printf("%v:%v\n", key, value)
    }
    //fmt.Println(flattened)

}

func decodeJSON(encodedJSON []byte) map[string]interface{} {

    decoded := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

func main() {
	load("working/ex3.json")
	
}