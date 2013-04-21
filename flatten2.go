package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"
	"strconv"
)

var lkey = ""

func flatten(inputJSON map[string]interface{}, lkey string, flattened *map[string]string) {
	for rkey, value := range inputJSON {
		key := lkey+rkey
		defineType(value, key, flattened)
		if _, ok := value.(string); ok {
			(*flattened)[key] = value.(string)
		} else if _, ok := value.([]interface{}); ok {
			// value is a list
			// flatten the values in the list
			for i := 0; i<len(value.([]interface{})); i++ {
				if _, ok := value.([]string); ok {
					stringI := string(i)
					(*flattened)[stringI] = value.(string)
					/// think this is wrong

				} else {
				flatten(value.([]interface{})[i].(map[string]interface{}), key+":"+strconv.Itoa(i)+":", flattened)
				}
			}
		} else {
			flatten(value.(map[string]interface{}), key+":", flattened)
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