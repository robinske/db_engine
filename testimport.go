package main

import (
	"io/ioutil"
	"db_engine/client"
	"log"
	"encoding/json"
	"fmt"
	//"github.com/likexian/simplejson"
	//"https://github.com/jmoiron/jsonq"
)

type JSON map[string]interface{}
type dictionary map[string]string
var flattened = dictionary {}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func decodeJSON(encodedJSON []byte) JSON {

	decoded := map[string]interface{} {}
	err := json.Unmarshal(encodedJSON, &decoded)
	check(err)

	return decoded
}


func main() {
	
	fileContents, err := ioutil.ReadFile("working/example.json")
	check(err)
	
	decoded := decodeJSON(fileContents)

	keys := []string{}

	for k := range decoded {
		keys = append(keys, k)
		fmt.Printf("%v: %v\n", k, decoded[k])
	}

	//fmt.Println(keys)
	//fmt.Printf("Value of key: %v\n", decoded["batters"])
	
	client.Connect(fileContents)

}