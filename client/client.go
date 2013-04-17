// CLIENT DIALS / WRITES

package client

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
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf[:]) // calls the read method on io reader variable r, sets instance to n
        if err != nil {
          return
        }
        fmt.Printf("Value: %s\n", string(buf[0:n])) // how to make this dynamic based on input?
        // can do println to add a buffer space between inputs
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
        // arbitrary input from another go program. of some kind.

        // rolodex of hackbright students
        // want this program to be able to save this information to the database
        // how to send information between two 

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