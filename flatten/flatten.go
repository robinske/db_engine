package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"
	"strconv"
)

var lkey = ""
var flattened = make(map[string]interface{})


func insert(inputJSON map[string]interface{}, oldkey string, flattened *map[string]interface{}) {
	for newkey, value := range inputJSON {
		key := oldkey+newkey
		if _, ok := value.(string); ok {
			(*flattened)[key] = value.(string)
		} else if _, ok := value.(float64); ok {
			(*flattened)[key] = value.(float64)
		} else if _, ok := value.(bool); ok {
			(*flattened)[key] = value.(bool)
		} else if _, ok := value.([]interface{}); ok {
			for i := 0; i<len(value.([]interface{})); i++ {
				if _, ok := value.([]string); ok {
					stringI := string(i)
					(*flattened)[stringI] = value.(string)
					/// think this is wrong
				} else if _, ok := value.([]int); ok {
					stringI := string(i)
					(*flattened)[stringI] = value.(int)
				} else {
					flatten(value.([]interface{})[i].(map[string]interface{}), key+":"+strconv.Itoa(i)+":", flattened)
				}
			}
		} else {
			flatten(value.(map[string]interface{}), key+":", flattened)
		}
	}
}

func load(filename string) []byte {
    encodedJSON, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    //fmt.Println(mappedJSON)
    return encodedJSON
}

func decodeJSON(encodedJSON []byte) map[string]interface{} {
    decodedJSON := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decodedJSON)
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Println(decoded)

    fmt.Printf("%v\n", decodedJSON["batters"].(map[string]interface{})["batter"].([]interface{})[0].(map[string]interface{})["id"])

    return decodedJSON
}

func main() {
	encodedJSON := load("example.json")
	decodedJSON := decodeJSON(encodedJSON)
	
	flatten(decodedJSON, lkey, &flattened)
    // for key, value := range flattened {
    // 	fmt.Printf("%v:%v\n", key, value)
    // }	
}