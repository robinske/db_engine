// SERVER LISTENS
// http://stackoverflow.com/questions/2886719/unix-sockets-in-go

package main

import (
    "net"
    "log"
    "io"
    "os"
    "strings"
    "fmt"
)

type structDict struct {
    Key string
    Value string
}

type cacheData map[string]string
var dictionary = cacheData {} // Declare global variable so not to overwrite

func echoServer(c net.Conn) {
    for {
        buf := make([]byte, 512) // makes a buffer to keep chunks that are read/written
        nr, err := c.Read(buf)
        if err == io.EOF {
            return
        }

        data := buf[0:nr]

        message := string(data)
        instruct, key, value := parseRequest(message)

        talkToDictionary(c, instruct, key, dictionary, value)

        fmt.Printf("Server received: %s", string(data))
        // _, err = c.Write(data) // this writes things for the client
        // if err != nil {
        //     log.Fatal(err)
        // }
    }
}

func parseRequest(message string) (instruct, key, value string) {
    msgSplit := strings.Split(message, " ")

    if len(msgSplit) == 0 {
        return
    }

    instruct = strings.TrimSpace(msgSplit[0])

    if len (msgSplit) == 1 {
        return
    }

    key = strings.TrimSpace(msgSplit[1])

    if len(msgSplit) == 2 {
        return
    }

    value = strings.TrimSpace(strings.Join(msgSplit[2:], " "))

    return
}

func talkToDictionary(c net.Conn, instruct, key string, dictionary cacheData, optionalValue...string) (cacheData) {

    value := strings.Join(optionalValue[:], " ")

    switch instruct {
        case "GET": get(c, key, dictionary)
        case "PUT": put(c, key, value, dictionary)
        //case "SAVE": save(key, value, instruct, dictionary)
        default: fmt.Println("try again idiot")
    }

    return dictionary
}

func get(c net.Conn, key string, dictionary cacheData) (value string) {

    value = dictionary[key]
    
    //fmt.Printf("Printing value: %s\n", value)
    byteValue := []byte(value)
    c.Write(byteValue) // sends the value back over to the client

    return
}

func put(c net.Conn, key, value string, dictionary cacheData) {

    dictionary[key] = value
    fmt.Println(dictionary)

    // ONCE THE DICTIONARY IS STRING/JSON - SEND IT OVER
    // byteDict := []byte(dictionary)
    // c.Write(byteDict)

    // ADD IF STATEMENT TO NOT OVERWRITE - NEW FUNCTION UPDATE WILL DO THAT
}

func save(key, value, instruct string, dictionary cacheData) {
    
    // WRITE TO DATABASE

    fo, err := os.OpenFile("outputs/output", os.O_RDWR|os.O_APPEND, 0666)
    fo.Seek(0,2) // 2 means go to the end of the file, 0 is the relative position to the end
    if err != nil {
        log.Fatal(err)
    }
    
    defer fo.Close()

    // _, err = fo.Write(dictionary)

    // WRITE TO DATABASE END
}

func main() {
    l, err := net.Listen("tcp", ":4127") // sets a listener, l, to port 4127
    if err != nil {
        log.Fatal(err)
        return
    }

    for {
        fd, err := l.Accept()
        if err != nil {
            log.Fatal(err)
            return
        }

        go echoServer(fd)
    }
}