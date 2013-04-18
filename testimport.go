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

func decodeJSON(encodedJSON []byte) {

	u := map[string]interface{} {}
	err := json.Unmarshal(encodedJSON, &u)
	check(err)

	keys := []string{}

	for k := range u {
		keys = append(keys, k)
		fmt.Printf("%v: %v\n", k, u[k])
	}

	fmt.Println(keys)
	fmt.Printf("Value of key: %v\n", u["batters"])

}


func main() {
	
	fileContents, err := ioutil.ReadFile("working/example.json")
	check(err)
	
	decodeJSON(fileContents)
	
	client.Connect(fileContents)

}