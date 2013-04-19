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

const (
    PORT = ":4127"
)

func reader(reader io.Reader) {
    buf := make([]byte, 1024)
    for {
        inputEnd, err := reader.Read(buf[:])
        if err != nil {
          return
        }
        fmt.Printf("%s\n", string(buf[0:inputEnd]))
    }
}

func main() {
    connection, err := net.Dial("tcp", PORT) 
    if err != nil {
        log.Fatal(err)
    }

    defer connection.Close()

    go reader(connection)
    
    for {
        input := bufio.NewReader(os.Stdin)
        rawMessage, err := input.ReadString('\n')
        if err != nil {
          log.Fatal(err)
        }

        message := strings.ToUpper(rawMessage)          // normalize input

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