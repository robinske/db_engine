// CLIENT DIALS / WRITES

package client

import (
    "net"
    "io"
    "log"
    "strings"
    "fmt"
)

func reader(reader io.Reader) {
    buf := make([]byte, 1024)
    for {
        inputEnd, err := reader.Read(buf[:]) // calls the read method on io reader variable r, sets instance to n
        if err != nil {
          return
        }
        fmt.Printf("Value: %s\n", string(buf[0:inputEnd])) // how to make this dynamic based on input?
        // can do println to add a buffer space between inputs
    }
}

func Connect(data string) {
    connection, err := net.Dial("tcp", ":4127") // sets a connection, c, to the port 4127
    if err != nil {
        log.Fatal(err)
    }

    defer connection.Close()

    go reader(connection)
    
    for {

        message := strings.ToUpper(data) // capitalize/normalize input

        if err != nil {
          log.Fatal(err)
        }

        if message != "" {
            _,err := connection.Write([]byte(message))

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