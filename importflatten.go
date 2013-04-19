package main

import (
	"fmt"
	//"db_engine/flat"
	"io/ioutil"
	"log"
	//"strings"
	"encoding/json"
)

type Dictionary map[string]string
type JSON map[string]interface{}

var lkey = ""

var flattenedJSON = Dictionary {}

func flatten(inputJSON JSON, lkey string, flattened *map[string]string) {
	for rkey, value := range inputJSON {
		key := lkey+rkey
		if _, ok := value.(string); ok {
			(*flattened)[key] = value.(string)
		} else {
			flatten(value.(map[string]interface{}), key+string(0), flattened)
			// 0 won't print anything but does contain a byte
			// maybe for testing use a unicode character

			
			// for key, value := range flattened {
			// 	flattenedJSON[key] = value.(string)
			// }
		}
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

    fmt.Println(flattened)

}

func decodeJSON(encodedJSON []byte) JSON {

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