// CLIENT DIALS / WRITES

package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
)

func main() {
  	resp, err := http.Get("http://localhost:4127")
  	if err == nil {
      fmt.Printf("errorrrrr")
    }
    if err != nil {
      log.Fatal(err)
  	}
  	body, err := ioutil.ReadAll(resp.Body)
  	resp.Body.Close()
  	if err != nil {
  		log.Fatal(err)
  	}
  	fmt.Printf("%s", body)
}