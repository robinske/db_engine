package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"
)

var lkey = ""
var nested = make(map[string]interface{})


// func insert(inputJSON map[string]interface{}, key string, nested *map[string]interface{}) {
// 	for key, value := range inputJSON {
// 		if _, ok := value.(string); ok {
// 			(*nested)[key] = value.(string)
// 		} else if _, ok := value.(float64); ok {
// 			(*nested)[key] = value.(float64)
// 		} else if _, ok := value.(bool); ok {
// 			(*nested)[key] = value.(bool)
// 		} else if _, ok := value.([]int); ok {
// 			(*nested)[key] = value.([]int)
// 		} else if _, ok := value.([]interface{}); ok {
// 			for index, value := range value.([]interface{}) {
// 				if _, ok := value.([]string); ok {
// 					stringI := string(index)
// 					(*nested)[stringI] = value.(string)
// 				} else if _, ok := value.([]int); ok {
// 					stringI := string(index)
// 					(*nested)[stringI] = value.(int)
// 				} else {
// 					insert(value.([]interface{})[index].(map[string]interface{}), key, nested)
// 				}
// 			}
// 		} else {
// 			insert(value.(map[string]interface{}), key, nested)
// 		}
// 	}
// }

func insert(inputJSON map[string]interface{}, nested *map[string]interface{}) {
	for key, value := range inputJSON {
		switch v := value.(type) {
		    case float64: (*nested)[key] = v
		    case string: (*nested)[key] = v
		    case bool: (*nested)[key] = v
		    case []interface{}: 
		    	//list := []interface{}{}
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

    //fmt.Printf("%v\n", decodedJSON["batters"].(map[string]interface{})["batter"].([]interface{})[0].(map[string]interface{})["id"])

    return decodedJSON
}

func main() {
	encodedJSON := load("example.json")
	decodedJSON := decodeJSON(encodedJSON)
	
	insert(decodedJSON, &nested)
    for key, value := range nested {
    	fmt.Printf("%v:%v\n", key, value)
    }	
}