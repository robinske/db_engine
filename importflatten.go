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

func flatten(inputJSON JSON, lkey string) JSON {
	for rkey, value := range inputJSON {
		key := lkey+rkey
		if _, ok := value.(map[string]string); ok {
			flattenedJSON[key] = value.(string)
		} else {
			flattened := flatten(inputJSON, lkey+"_")
			for key, value := range flattened {
				flattenedJSON[key] = value.(string)
			}
		}
	}
	return flattened
}

func load(filename string) {
    fileContents, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    mappedJSON := decodeJSON(fileContents)
    flatten(mappedJSON, lkey)

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
	fmt.Println(flattenedJSON)
}