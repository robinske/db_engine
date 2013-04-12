// CLIENT DIALS / WRITES

package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  	http.HandleFunc("/", handler)
  	resp, err := http.Get("http://localhost:4127")
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