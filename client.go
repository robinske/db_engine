// CLIENT DIALS / WRITES

package main

import (
    "net"
    "log"
    "strings"
    "bufio"
    "os"
    "fmt"
)

var DATABASE string

const (
    PORT = ":4127"
)

func main() {

    connection, err := net.Dial("tcp", PORT) 
    if err != nil {
        log.Fatal(err)
    }

    input := bufio.NewReader(os.Stdin)
    buf := make([]byte, 1024)

    defer connection.Close()

    if len(os.Args) == 2 {
        DATABASE = strings.TrimSpace(os.Args[1])

        _, err := connection.Write([]byte("DATABASE:> "+DATABASE+"\n")) // have to hit enter again??? MIGHT BE BECAUSE OF LINE 51  
        if err != nil {
            log.Fatal(err)
        }

    } else {
        connection.Write([]byte("Please load a database"))
    }

    inputEnd, err := connection.Read(buf[:])
    if err != nil {
        return
    }

    fmt.Printf("%s\n", string(buf[0:inputEnd]))

    for {
        fmt.Printf(">>> ")
        rawMessage, err := input.ReadString('\n')

        if err != nil {
          log.Fatal(err)
        }

        message := strings.TrimSpace(strings.ToUpper(rawMessage))          // normalize input

        if message != "" {

            connection.Write([]byte(message))

            inputEnd, err := connection.Read(buf[:])
            if err != nil {
                return
            }

            fmt.Printf("%s\n", string(buf[0:inputEnd]))
        }

        if message == "QUIT" {
            fmt.Println("Goodbye!")
            // connection.Write([]byte(strings.TrimSpace(session)+" has been disconnected\n"))   // name the session
            connection.Close()
            return
        }
    }
}