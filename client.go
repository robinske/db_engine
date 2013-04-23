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

const (
    PORT = ":4127"
)

func main() {
    connection, err := net.Dial("tcp", PORT) 
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Please Enter Your Session Name>> ")
    sess := bufio.NewReader(os.Stdin)
    session, err := sess.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }

    defer connection.Close()

    input := bufio.NewReader(os.Stdin)

    buf := make([]byte, 1024)
       
    for {
        fmt.Printf(">>> ")
        rawMessage, err := input.ReadString('\n')

        if err != nil {
          log.Fatal(err)
        }

        message := strings.ToUpper(rawMessage)          // normalize input

        if strings.TrimSpace(message) == "QUIT" {
            fmt.Println("Goodbye!")
            connection.Write([]byte(strings.TrimSpace(session)+" has been disconnected\n"))   // name the session
            connection.Close()
            return
        }

        if message != "" {

            _,err := connection.Write([]byte(message))
            if err != nil {
              log.Fatal(err)
              break
            }
            
            inputEnd, err := connection.Read(buf[:])
            if err != nil {
                return
            }

            fmt.Printf("%s\n", string(buf[0:inputEnd]))
        }
    }
}