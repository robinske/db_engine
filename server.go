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

// MAKE SURE EACH FUNCTION ONLY DOES ONE THING

type dictionary map[string]string
var cacheData = dictionary {} // Declare global variable so not to overwrite

// var connection net.Conn getting error but could try setting this as a global

const (
    PORT = ":4127"
)

func echoServer(connection net.Conn) (data []byte) {
    for {
        buf := make([]byte, 512) // makes a buffer to keep chunks that are read/written
        inputEnd, err := connection.Read(buf) // COMMENT THIS BETTER
        if err == io.EOF {
            return
        }

        data = buf[0:inputEnd]

        message := string(data)
        instruction, key, value := parseRequest(message) // how to take these out of the function?

        talkToDictionary(connection, instruction, key, value) // how to take these out of the function?

        fmt.Printf("Server received: %s", string(data))

    }

    return
}

func parseRequest(message string) (instruction, key, value string) {
    msgSplit := strings.Split(message, " ")

    if len(msgSplit) == 0 {
        return
    }

    instruction = strings.TrimSpace(msgSplit[0])

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

func talkToDictionary(connection net.Conn, instruct, key string, optionalValue...string) {

    value := strings.Join(optionalValue[:], " ")

    switch instruct {
        case "GET": get(connection, key)
        case "PUT": put(connection, key, value)
        //case "SAVE": save(key, value, instruct, dictionary)
        default: fmt.Println("try again idiot")
    }

}

func get(connection net.Conn, key string) (value string) {

    value = cacheData[key]
    
    //fmt.Printf("Printing value: %s\n", value)
    byteValue := []byte(value)
    connection.Write(byteValue) // sends the value back over to the client

    return
}

func put(connection net.Conn, key, value string) {

    // make clear for which dictionary for when multiple clients are dealing with different cache

    cacheData[key] = value
    fmt.Println(cacheData)
    // Give the client confirmation that this worked

    // ONCE THE DICTIONARY IS STRING/JSON - SEND IT OVER
    // byteDict := []byte(dictionary)
    // c.Write(byteDict)

    // ADD IF STATEMENT TO NOT OVERWRITE - NEW FUNCTION UPDATE WILL DO THAT
}

func save(key, value, instruct string) {
    
    // WRITE TO DATABASE

    fo, err := os.OpenFile("outputs/output", os.O_RDWR|os.O_APPEND, 0666) // open file outside of this function
    fo.Seek(0,2) // 2 means go to the end of the file, 0 is the relative position to the end
    if err != nil {
        log.Fatal(err)
    }
    
    defer fo.Close()

    // _, err = fo.Write(dictionary)

    // WRITE TO DATABASE END
}

func main() {
    listener, err := net.Listen("tcp", PORT)
    if err != nil {
        log.Fatal(err)
        return
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
            return
        }

        go echoServer(conn)
    }
}