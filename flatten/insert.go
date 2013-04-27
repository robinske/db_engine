package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"
)

var lkey = ""
var nested = make(map[string]interface{})


func insert(inputJSON map[string]interface{}, nested *map[string]interface{}) {
	for key, value := range inputJSON {
		switch v := value.(type) {
		    case float64: (*nested)[key] = v
		    case string: (*nested)[key] = v
		    case bool: (*nested)[key] = v
		    case []interface{}: 
		    	for _, val := range v {
		    		if _, ok := val.(string); ok {
						v = append(v, val.(string))
						(*nested)[key] = v
					} else if _, ok := val.(int); ok {
						v = append(v, val.(string))
						(*nested)[key] = v
					} else {
						insert(val.(map[string]interface{}), nested)
						(*nested)[key] = v
					}
				}
			case map[string]interface{}: 
				(*nested)[key] = v
				insert(v, nested)
		    default:
		        fmt.Println("ABLAKJFALSKDFJ WHAT DID I MISS", v)
		}
	}
}


func load(filename string) []byte {
    encodedJSON, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    return encodedJSON
}

func decodeJSON(encodedJSON []byte) map[string]interface{} {
    decodedJSON := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decodedJSON)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%v\n", decodedJSON["batters"].(map[string]interface{})["batter"].([]interface{})[0].(map[string]interface{})["id"])

    return decodedJSON
}

func main() {
	encodedJSON := load("example.json")
	decodedJSON := decodeJSON(encodedJSON)
	
	insert(decodedJSON, &nested)

    // for key, value := range nested {
    // 	fmt.Printf("%v:%v\n", key, value)
    // }	

    fmt.Printf("%v\n", nested["batters"]["batter"])
}