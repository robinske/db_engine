// CLIENT DIALS / WRITES

package main

  import (
    "net"
    "io"
    "log"
    "strings"
    "bufio"
    "os"
    "fmt"
  )

func reader(r io.Reader) {
  buf := make([]byte, 1024) // makes a list of bytes / why 1024? why a list? whyyyyy
  for {
    n, err := r.Read(buf[:]) // calls the read method on r (which is passed in but what is r?? sets this instance to n
    if err != nil {
      return
    }
    fmt.Printf("Client dialed: %s", string(buf[0:n])) // can do println to add a buffer space between inputs
  }
}

func displayResults() {

  
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
    rawMessage, err := input.ReadString('\n') // message resets with each new input
    message := strings.ToUpper(rawMessage) // capitalize/normalize input

    if err != nil {
      log.Fatal(err)
    }
    if message != "" {
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