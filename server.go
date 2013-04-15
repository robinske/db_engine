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

type cacheData map[string]string
var dictionary = cacheData {}

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

        talkToDictionary(instruct, key, value, dictionary)

        println("Server received:", string(data))
        _, err = c.Write(data)
        if err != nil {
            log.Fatal(err)

        }
    }
}

func parseRequest(message string) (string, string, string) {
    msgSplit := strings.Split(message, " ")

    switch len(msgSplit) {
        case 1:
        case 2:
        case 3:
        default: // enter/no length
    }
    
    if len(msgSplit) <= 2 {             // handle for get requests
        return instruct, key, "nil"     // this might be why it's breaking??
    }

    instruct := msgSplit[0]
    key := msgSplit[1]
    value := strings.Join(msgSplit[2:], " ")

    fmt.Printf("instruct: %s, key: %s, value: %s", instruct, key, value)
    // need to handle message length error - only works for 3+ word inputs right now

    return instruct, key, value
}

func talkToDictionary(instruct, key, value string, dictionary cacheData) (cacheData) {

    switch instruct {
        case "GET": get(key, dictionary)
        case "PUT": put(key, value, dictionary)
        case "SAVE": save(key, value, instruct, dictionary)
        default: fmt.Println("try again idiot")
    }

    if instruct == "GET" {
        get(key, dictionary)
    }

    getValue := get(key, dictionary)
    println("value is always called:", getValue) // this works when put is called!!!

    fmt.Println(dictionary)
    return dictionary
}

func get(key string, dictionary cacheData) (string) {

    getKey := key

    //fmt.Printf("We'll get the value of key %s\n", key)
    value := dictionary[getKey]
    fmt.Printf("Printing value: %s end of value\n", value)

    // RETURN DATA TO THE CLIENT
    return value
}

func put(key, value string, dictionary cacheData) {
    // fmt.Printf("You called the %s function\n", instruct)
    putKey := key
    putValue := value

    fmt.Printf("We'll put %s:%s in the dictionary\n", putKey, putValue)
    dictionary[putKey] = putValue
    //fmt.Println(data[key]) // should print value
}

func save(key, value, instruct string, dictionary cacheData) {
    // want to save dictionary to disk on exit or when explicitly called

    // save the operations that wrote to the database
    // JUST THE THINGS THAT MODIFY THE DATABASE
    // ASYNCRONOUSLY SAVING
    // QUEUE UP DISK WRITES/ASSISTANT FUNCTION (WHEN YOU HIT SAVE THE QUEUE WILL BE WRITTEN TO DISK)
    
    // WRITE TO DATABASE

    fo, err := os.OpenFile("outputs/output", os.O_RDWR|os.O_APPEND, 0666)
    fo.Seek(0,2) // 2 means go to the end of the file, 0 is the relative position to the end
    if err != nil {
        log.Fatal(err)
    }
    
    defer fo.Close()

    //_, err = fo.Write(dictionary)

    // WRITE TO DATABASE END
}

func main() {
    l, err := net.Listen("tcp", ":4127")
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