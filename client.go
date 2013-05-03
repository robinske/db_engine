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
    BUFFER_SIZE = 1024
)

func main() {

    connection, err := net.Dial("tcp", PORT) 
    if err != nil {
        log.Fatal(err)
    }

    input := bufio.NewReader(os.Stdin)
    buf := make([]byte, BUFFER_SIZE)

    defer connection.Close()

    if len(os.Args) == 2 {
        DATABASE = strings.TrimSpace(os.Args[1])

        _, err := connection.Write([]byte("DATABASE:> "+DATABASE))
        if err != nil {
            log.Fatal(err)
        }
        inputEnd, err := connection.Read(buf[:])
        if err != nil {
            return
        }
        fmt.Printf("%s\n", string(buf[0:inputEnd]))
        
    } else {
        fmt.Println("Please load a database")
    }

    for {
        fmt.Printf(">>> ")
        rawMessage, err := input.ReadString('\n')

        if err != nil {
            log.Fatal(err)
        }

        message := strings.TrimSpace(strings.ToUpper(rawMessage))          // normalize input

        if message != "" {

            connection.Write([]byte(message))

            for {
                inputEnd, err := connection.Read(buf[:])
                fmt.Printf("%s\n", string(buf[0:inputEnd]))
                if inputEnd < BUFFER_SIZE {
                    break
                }
                if err != nil {
                    return
                }            
            }
        }

        if message == "QUIT" {
            connection.Close()
            return
        }
    }
}