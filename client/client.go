// CLIENT DIALS / WRITES

package client

import (
    "net"
    "io"
    "log"
    "strings"
    "fmt"
    "bufio"
    "os"
)

func reader(reader io.Reader) {
    buf := make([]byte, 10000)
    for {
        inputEnd, err := reader.Read(buf[:]) // calls the read method on io reader variable r, sets instance to n
        if err != nil {
          return
        }
        fmt.Printf("Value: %s\n", string(buf[0:inputEnd])) // how to make this dynamic based on input?
        // can do println to add a buffer space between inputs
    }
}

func Connect(data []byte) {
    connection, err := net.Dial("tcp", ":4127") // sets a connection, c, to the port 4127
    
    if err != nil {
        log.Fatal(err)
    }

    defer connection.Close()

    go reader(connection)

    message := string(data) // capitalize/normalize input
    
    for {

        if err != nil {
          log.Fatal(err)
        }

        if message != "" {
            _, err := connection.Write([]byte(message))

            if err != nil && err != io.EOF {
              log.Fatal(err)
              break
            }
        }

        input := bufio.NewReader(os.Stdin)
        rawMessage, err := input.ReadString('\n')
        message = strings.ToUpper(rawMessage) // resets message to interaction

        if err != nil {
            log.Fatal(err)
            break
        }
    }
}