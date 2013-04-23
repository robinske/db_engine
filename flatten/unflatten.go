package main

import (
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
	"strings"
	"strconv"
)

type Dictionary map[string]interface{}

var cacheData = Dictionary {}
var lkey = ""
var flattened = make(map[string]interface{})
var jsonString = `{`

func load(filename string) {
    fileContents, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    mappedJSON := decodeJSON(fileContents)

    flatten(mappedJSON, lkey, &flattened)

    for key, value := range flattened {
        fmt.Printf("%v:%v\n", key, value)
    }
    
    for k,v := range flattened {
        k = strings.ToUpper(k)
        if _, ok := v.(string); ok {
            v = strings.ToUpper(v.(string))
            cacheData[k] = v.(string)
        } else if _, ok := v.(float64); ok {
            v = v.(float64)
            cacheData[k] = v.(float64)
        } else if _, ok := v.(bool); ok {
            v = v.(bool)
            cacheData[k] = v.(bool)
        } else {
            fmt.Println("JSON file format error")
        }
    }
}

func decodeJSON(encodedJSON []byte) map[string]interface{} {
    decoded := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

func flatten(inputJSON map[string]interface{}, lkey string, flattened *map[string]interface{}) {
    for rkey, value := range inputJSON {
        key := lkey+rkey
        if _, ok := value.(string); ok {
            (*flattened)[key] = value.(string)
        } else if _, ok := value.(float64); ok {
            (*flattened)[key] = value.(float64)
        } else if _, ok := value.(bool); ok {
            (*flattened)[key] = value.(bool)
        } else if _, ok := value.([]float64); ok { // type check for a list of integers not working - is this valid JSON though?
            (*flattened)[key] = value.([]float64)
        } else if _, ok := value.([]interface{}); ok {
            for i := 0; i<len(value.([]interface{})); i++ {
                if _, ok := value.([]string); ok {
                    stringIndex := string(i)
                    (*flattened)[stringIndex] = value.(string)
                } else {
                    flatten(value.([]interface{})[i].(map[string]interface{}), 
                            key+":"+strconv.Itoa(i)+":", flattened)
                }
            }
        } else {
            flatten(value.(map[string]interface{}), key+":", flattened)
        }
    }
}

var uniqueKeys = map[string]string {}


func unflatten(flattened map[string]interface{}, jsonString string) string {
	
	for k,v := range flattened {
		if strings.Contains(k, ":") == false {
			if _, ok := v.(string); ok {
				jsonString = jsonString+"\""+k+"\":\""+v.(string)+"\","
			} else if _, ok := v.(float64); ok {
				jsonString = jsonString+"\""+k+"\":"+strconv.FormatFloat(v.(float64), 'f', -1, 64)+","
			} else if _, ok := v.(bool); ok {
				jsonString = jsonString+"\""+k+"\":"+strconv.FormatBool(v.(bool))+","
			} else {
				fmt.Println("da fuck happened?")
			}
		} else {
			splitKey := strings.SplitN(k, ":", 2)
			newKey := splitKey[0]
			newValue := strings.Join(splitKey[1:], " ")

			// listIndex, _ := strconv.Atoi(newKey)
			listElems := []string{}
			listElems = append(listElems, newValue)
			
			fmt.Println(listElems)
			uniqueKeys[newKey] = ""

			innerCache := map[string]interface{}{}
			innerCache[newValue] = v
			unflatten(innerCache, jsonString)
			fmt.Println("newKey: ", newKey, "newValue: ", newValue)
		}
	}
	for key := range uniqueKeys {
		jsonString = jsonString+"\""+key+"\":NESTED VALUES HERE,"
	}
	jsonString = jsonString+"}"
	return jsonString

}

func encode(sjson string) {
	b, err := json.Marshal(sjson)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b)
	fmt.Println(string(b))
}

func main() {
	
	load("example.json")
	sjson := unflatten(cacheData, jsonString)

	fmt.Println(sjson)
	//fmt.Println(jsonString)
	// sjson := `{"Name":"Alice", "Body":"Hello", "Time":{"Day":"Monday", "Hour":"Afternoon"}}`

	//encode(sjson)
	
}