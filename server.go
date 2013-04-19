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
    // "sync"
)

// MAKE SURE EACH FUNCTION ONLY DOES ONE THING



type Dictionary map[string]string
type Collection map[string]map[string]string
type JSON map[string]interface{}

// var locker = struct{
//     sync.RWMutex
//     cacheData map[string]int
// }{cacheData := make(map[string]int)}

var cacheData = Dictionary {} // Declare global variable so not to overwrite - HOW TO IMPLEMENT THIS FOR MULTIPLE DICTIONARIES/COLLECTIONS?
var queue []byte // what will be written to disk

const (
    PORT = ":4127"
)

func echoServer(connection net.Conn) {      // this function does too many things. need to separate it
    for {
        buf := make([]byte, 10000) // make buffer size infinite/flexible based on data input?
        inputEnd, err := connection.Read(buf)
        if err == io.EOF {
            return
        }

        data := buf[0:inputEnd]
        message := string(data)
        instruction, key, value := parseRequest(message) 

        callCacheData(connection, instruction, key, value)

        fmt.Printf("Server received: %s", message)
    }
}

func parseRequest(message string) (instruction, key, value string) {
    
    msgSplit := strings.Split(message, " ")
    
    if len(msgSplit) == 0 { return }
    instruction = strings.TrimSpace(msgSplit[0])

    if len(msgSplit) == 1 { return }
    key = strings.TrimSpace(msgSplit[1])

    if len(msgSplit) == 2 { return }
    value = strings.TrimSpace(strings.Join(msgSplit[2:], " "))

    return
}

func callCacheData(connection net.Conn, instruction, key string, optionalValue...string) {

    value := strings.Join(optionalValue[:], " ")

    // keyExists := check(key)
    // TEST FOR KEY EXISTENCE HERE, call functions differently based on value?

    switch instruction {
        case "CREATE": {
            var key = Collection {}
            collection := key
            fmt.Println(collection)
            create(connection, collection)
        }
        case "GET": get(connection, key)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "LOAD": {
            filename := key
            load(connection, filename)
        }
        case "SHOW": show(connection, key)
        case "QUIT": quit(connection)
        case "REMOVE": remove(connection, key)
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

func trackCollections() {
}

func check(key string) (exists bool) {
    // value, ok := cacheData[key]         // check if key is valid
    // if ok {
    //         byteValue := []byte(value)
    //         connection.Write(byteValue)
    // } else {
    //         connection.Write([]byte("key not found"))
    // }
    return
}

func quit(connection net.Conn) {
    // write entire dictionary to disk here
    connection.Write([]byte("Connection has been terminated"))
    connection.Close()
}
func create(connection net.Conn, collection Collection) {
    // prefix 

}

func get(connection net.Conn, key string) (value string) {

    value, ok := cacheData[key]         // check if key is valid
    if ok {
            byteValue := []byte(value)
            connection.Write(byteValue)
    } else {
            connection.Write([]byte("key not found"))
    }
    return  
}

func set(connection net.Conn, key, value string) {

    // make clear for which dictionary for when multiple clients are dealing with different cache
    // ADD IF STATEMENT TO NOT OVERWRITE - NEW FUNCTION UPDATE WILL DO THAT
    _, ok := cacheData[key]         // check if key is valid
    if ok {
        connection.Write([]byte(key+" already added. To modify, UPDATE key"))
    } else {
        cacheData[key] = value      
        connection.Write([]byte("Added "+key+":"+value))
    }
}

func update(connection net.Conn, key, value string) {
    
    _, ok := cacheData[key]         // check if key is valid
    if ok {
        cacheData[key] = value      // overwrite
        connection.Write([]byte("Updated "+key+":"+value))
    } else {
        connection.Write([]byte(key+" not yet added. Please set"))
    }
}

func load(connection net.Conn, filename string) {
    // ERROR HANDLE FILE NOT RECOGNIZED
    fileContents, err := ioutil.ReadFile(filename)
    connection.Write([]byte("Loaded "+filename+" to collection X"))
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
    // fmt.Println(cacheData)
    // show things in database
    // i.e. "show keys"
    switch key {
        case "KEYS": {
            keys := []string{}
            for k := range cacheData {
                keys = append(keys, k)
            }
            connection.Write([]byte(strings.Join(keys, ", ")))
        }
        // case "VALUES": {
        //     for k, v := range cacheData {
        //         connection.Write([]byte("Key: "+k+", Value: "+v))
        //     }
        //     connection.Write([]byte(strings.Join(keys, ", ")))
        // }
        case "COLLECTIONS": {
            // SHOW THE DIFFERENT DICTIONARIES IN CACHE
        }
        default: connection.Write([]byte("Invalid request"))

    }
}

func remove(connection net.Conn, key string) {
    
    _, ok := cacheData[key]         // check if key is valid
    if ok {
        delete(cacheData, key)
        connection.Write([]byte(key+" has been removed"))
    } else {
        connection.Write([]byte("No key: "+key))
    }  
}

func openDisk() *os.File {
    disk, err := os.OpenFile("outputs/output", os.O_RDWR|os.O_APPEND, 0666) // open file outside of this function
    if err != nil {
        log.Fatal(err)
    }
   
    defer disk.Close()
    return disk
}

// func queueWrites() {
//     // use some global variable byte string to queue up stuff
//     // update the byte string with things to save
// }

func save(disk *os.File) {
    END := 2
    disk.Seek(0,END)
    _, err := disk.Write(queue)
        if err != nil {
        log.Fatal(err)
    }
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
        // disk := openDisk()
        // go save(disk)
    }
}