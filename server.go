// SERVER LISTENS
// http://stackoverflow.com/questions/2886719/unix-sockets-in-go

package main

import (
    "net"
    "log"
    "io"
    "io/ioutil"
    "os"
    "strings"
    "fmt"
    "encoding/json"
    // "db_engine/client"
)

// MAKE SURE EACH FUNCTION ONLY DOES ONE THING

type dictionary map[string]string
type JSON map[string]interface{}

var cacheData = dictionary {} // Declare global variable so not to overwrite

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
        case "LOAD": {
            filename := key
            load(connection, filename)
        }
        case "SHOW": show(connection, key)
        //case "SAVE": save(key, value, instruct, dictionary)
        default: fmt.Println("try again idiot")
    }
}

func get(connection net.Conn, key string) (value string) {

    //value = cacheData[key]

// CHECK IF KEY IS IN DICTIONARY

    value, ok := cacheData[key]
    if ok {
            byteValue := []byte(value)
            connection.Write(byteValue) // sends the value back over to the client
    } else {
            connection.Write([]byte("key not found"))
    }
    return   
}

func put(connection net.Conn, key, value string) {

    // make clear for which dictionary for when multiple clients are dealing with different cache

    cacheData[key] = value
    fmt.Println(cacheData)
    msg := "Added "+key+":"+value
    connection.Write([]byte(msg))
    // Give the client confirmation that this worked

    // ADD IF STATEMENT TO NOT OVERWRITE - NEW FUNCTION UPDATE WILL DO THAT
}

func load(connection net.Conn, filename string) {
    fileContents, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }
    mappedJSON := decodeJSON(fileContents)
    for k,v := range mappedJSON {                   // RIGHT HERE IS WHERE THERE ARE ISSUES. 
        k = strings.ToUpper(k)
        v = strings.ToUpper(v.(string))             // NEED TO EITHER DO A RECURSIVE SWITCH OR
                                                    // FLATTEN THE KEYS
        cacheData[k] = v.(string)
    }
    fmt.Printf("%v", mappedJSON)
}

func decodeJSON(encodedJSON []byte) JSON {

    decoded := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

func show(connection net.Conn, key string) {
    fmt.Println(cacheData)
    // show things in database
    // i.e. "show keys"
    switch key {
        case "KEYS": {
            keys := []string{}
            for k := range cacheData {
                keys = append(keys, k)
                fmt.Printf("%v\n", k)
            }
            keystring := strings.Join(keys, ", ")
            connection.Write([]byte(keystring))
        }
        default: connection.Write([]byte("Invalid request"))

    }
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