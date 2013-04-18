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
type netConn net.Conn

var cacheData = dictionary {} // Declare global variable so not to overwrite
//var connection = netConn {}

// var connection net.Conn getting error but could try setting this as a global

const (
    PORT = ":4127"
)

func echoServer(connection net.Conn) (data []byte) {
    for {
        buf := make([]byte, 10000) // make buffer size infinite/flexible based on data input?
        inputEnd, err := connection.Read(buf)
        if err == io.EOF {
            return
        }

        data = buf[0:inputEnd]

        message := string(data)

        instruction, key, value := parseRequest(message) // how to take these out of the function?

        callCacheData(connection, instruction, key, value) // how to take these out of the function?

        fmt.Printf("Server received: %s", message)
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

func callCacheData(connection net.Conn, instruct, key string, optionalValue...string) {

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

// CHECK IF KEY IS IN DICTIONARY

    // dict := map[string]int {"foo" : 1, "bar" : 2}
    // value, ok := dict["baz"]
    // if ok {
    //         fmt.Println("value: ", value)
    // } else {
    //         fmt.Println("key not found")
    // }
    
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

func show(connection net.Conn) {

    // show things in database
    // i.e. "show keys"

}

func openDisk() {
    END := 2
    fo, err := os.OpenFile("outputs/output", os.O_RDWR|os.O_APPEND, 0666) // open file outside of this function
    if err != nil {
        log.Fatal(err)
    }
    fo.Seek(0,END) // 2 means go to the end of the file, 0 is the relative position to the end
    defer fo.Close()

    // _, err = fo.Write(dictionary)

}

func save(key, value, instruct string) {

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