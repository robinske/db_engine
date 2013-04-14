// CLIENT DIALS / WRITES

package main

  import (
    "net"
    "io"
    "log"
    "fmt"
    "strings"
    "bufio"
    "os"
  )

func reader(r io.Reader) {
  buf := make([]byte, 1024) // makes a list of bytes / why 1024? why a list? whyyyyy
  for {
    n, err := r.Read(buf[:]) // calls the read method on r (which is passed in but what is r?? sets this instance to n
    if err != nil {
      return
    }
    println("Client dialed:", string(buf[0:n]))
  }
}

func saveToDict(message string) {



  
}

func parseRequest(message string) {
  msgSplit := strings.Split(message, " ")

  fmt.Printf("%q\n", message)
  fmt.Printf("%q\n", msgSplit)

  if msgSplit[0] == "GET" {
    // GET function
    fmt.Printf("you called the %s function\n", msgSplit[0])
  }

  if msgSplit[0] == "PUT" {
    // PUT function
    fmt.Printf("you called the %s function\n", msgSplit[0])
  }

  if msgSplit[0] == "SAVE" {
    // WRITE TO FILE function
    fmt.Printf("you called the %s function\n", msgSplit[0])
  }
}

func main() {
  c, err := net.Dial("tcp", ":4127") // sets a connection, c, to the port 4127
  if err != nil {
    log.Fatal(err)
  }
  defer c.Close()

  go reader(c) // concurrent process / Goroutine
  
  for {
    input := bufio.NewReader(os.Stdin)
    message, err := input.ReadString('\n') // message resets with each new input
    if err != nil {
      log.Fatal(err)
    }
    if message != "" {
        parseRequest(message)
        _,err := c.Write([]byte(message))
        if err != nil {
          log.Fatal(err)
          break
        }
    }
    if err != nil {
        log.Fatal(err)
        break
    }
  }
}