// CLIENT DIALS / WRITES

package main

  import (
    "net"
    "io"
    "log"
    "fmt"
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

  func writer(w io.Writer) {

  }

  func main() {
    c, err := net.Dial("tcp", ":4127") // sets a connection, c, to the port 4127
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    go reader(c) // concurrent process / Goroutine
    
    var message string
    // insert a pointer? ">>"
    fmt.Scanf("%x", &message) // only accepts the first word with &message, accepts/concatenates words with just message)
    // string module might help with this

    for {
        if message != "" {
            _,err := c.Write([]byte(message)) // this would be the place to pass in data...  
            if err != nil {
              log.Fatal(err)
              break
            }
        }
        if err != nil {
            log.Fatal(err)
            break
        }
        fmt.Scanf("%s", message)
    }
  }