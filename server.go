
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
    "strconv"
    "sync"
    "sort"
)

// MAKE SURE EACH FUNCTION ONLY DOES ONE THING

type Dictionary map[string]interface{}
type Collection map[string]map[string]string
type JSON map[string]interface{}

var lock = struct{
    sync.RWMutex
    cacheData map[string]interface{}
}{cacheData: make(map[string]interface{})}

// other locking options have a single goroutine that's responsible for applying changes ot the database only one channel to send the mutations -- will manage a queue of requests

var flattened = make(map[string]interface{})
var lkey = ""
var jsonString = ""
var state = true
var DATABASE string

const (
    PORT = ":4127"
    END = 2
    LOGFILE = "outputs/log.txt"
    //DATABASE = "outputs/output"
)

func echoServer(connection net.Conn) {      // this function does too many things. need to separate it

    for {
        buf := make([]byte, 10000)          // use bytes library for this
        inputEnd, err := connection.Read(buf)
        if err == io.EOF {
            return
        }

        dataInput := buf[0:inputEnd]
        message := string(dataInput)

        instruction, key, value := parseRequest(message)
        if instruction == "SET" || instruction == "UPDATE" || instruction == "REMOVE" {
            saveLog(dataInput)
        }

        callCacheData(connection, instruction, key, value)

        fmt.Printf("Server received: %s", message)
    }
}

func quit(connection net.Conn) {

    save()

    // for SHUT DOWN - os.Exit() - after you save

    return
}

func save() {
    
    data := encode()
    disk := openDisk(DATABASE)
    defer disk.Close()

    disk.Seek(0,END)
    _, err := disk.Write([]byte(data))
    if err != nil {
        log.Fatal(err)
    }

    clearLog(LOGFILE)   // can clear the log once written to stable storage
}

func parseRequest(message string) (instruction, key, value string) {
    
    msgSplit := strings.Fields(message)
    
    if len(msgSplit) == 0 { return }
    instruction = msgSplit[0]

    if len(msgSplit) == 1 { return }
    key = msgSplit[1]

    if len(msgSplit) == 2 { return }
    value = strings.Join(msgSplit[2:], " ")

    return
}

func callCacheData(connection net.Conn, instruction, key string, optionalValue...string) (DATABASE string) {

    value := strings.Join(optionalValue[:], " ")

    switch instruction {
        // case "CREATE": create(connection, collection)
        case "DATABASE:>": setDB(connection, key)
        case "GET": get(connection, key)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "LOAD": {
            filename := key
            load(connection, filename)
        }
        case "SHOW": show(connection, key)
        case "REMOVE": remove(connection, key)
        case "QUIT": quit(connection)
        case "SAVE": save()
        case "CLEAR": clearLog(LOGFILE)
        default: connection.Write([]byte("Instruction not recognized"))
    }
    return
}

func create(connection net.Conn, collection Collection) {
    // prefix 
}

func get(connection net.Conn, key string) {

    value, ok := lock.cacheData[key]         // check if key is valid
    values := []string{}

    if ok {
            connection.Write([]byte(value.(string)))
    } else {
        for k := range lock.cacheData {      // NO LONGER HASHING, O(N)
            if strings.Contains(k, key) {
                lock.RLock()
                v := lock.cacheData[k]
                lock.RUnlock()
                values = append(values, k+": "+v.(string))
            }
        }
        if len(values) == 0 {
            connection.Write([]byte("No values found"))
        } else {
            connection.Write([]byte(strings.Join(values, " \n")))
        }
    }
}

func set(connection net.Conn, key, value string) {

    // make clear for which dictionary for when multiple clients are dealing with different cache
    _, ok := lock.cacheData[key]         // check if key is valid
    if ok {
        connection.Write([]byte(key+" already added. To modify, UPDATE key"))
    } else {
        lock.Lock()
        lock.cacheData[key] = value
        lock.Unlock()
        connection.Write([]byte("Added "+key+":"+value))
    }
}

func update(connection net.Conn, key, value string) {
    
    _, ok := lock.cacheData[key]         // check if key is valid
    if ok {
        lock.Lock()
        lock.cacheData[key] = value     // overwrite
        lock.Unlock()
        connection.Write([]byte("Updated "+key+":"+value))
    } else {
        connection.Write([]byte(key+" not yet added. Please set"))
    }
}

func load(connection net.Conn, filename string) {
    fileContents, err := ioutil.ReadFile(filename)      // need better error handling here -- if file does not exist don't break
    if err != nil {
        log.Fatal(err)
        return
    }
    
    mappedJSON := decodeJSON(fileContents)

    flatten(mappedJSON, lkey, &flattened)

    // for key, value := range flattened {
    //     fmt.Printf("%v:%v\n", key, value)
    // }
    
    for k,v := range flattened {
        k = strings.ToUpper(k)
        if _, ok := v.(string); ok {
            v = strings.ToUpper(v.(string))
            lock.cacheData[k] = v.(string)
        } else if _, ok := v.(float64); ok {
            v = v.(float64)
            lock.cacheData[k] = v.(float64)
        } else if _, ok := v.(bool); ok {
            v = v.(bool)
            lock.cacheData[k] = v.(bool)
        } else {
            fmt.Println("JSON file format error")
        }
    }
    connection.Write([]byte("Loaded "+filename+" to collection X"))
    return
}

func decodeJSON(encodedJSON []byte) map[string]interface{} {
    decoded := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

func flatten(inputJSON map[string]interface{}, lkey string, flattened *map[string]interface{}) {
    for rkey, value := range inputJSON {
        key := lkey+rkey
        if _, ok := value.(string); ok {
            (*flattened)[key] = value.(string)
        } else if _, ok := value.(float64); ok {
            (*flattened)[key] = value.(float64)
        } else if _, ok := value.(bool); ok {
            (*flattened)[key] = value.(bool)
        } else if _, ok := value.([]float64); ok { // type check for a list of integers not working - is this valid JSON though?
            (*flattened)[key] = value.([]float64)
        } else if _, ok := value.([]interface{}); ok {
            for i := 0; i<len(value.([]interface{})); i++ {
                if _, ok := value.([]string); ok {
                    stringIndex := string(i)
                    (*flattened)[stringIndex] = value.(string)
                } else {
                    flatten(value.([]interface{})[i].(map[string]interface{}), 
                            key+":"+strconv.Itoa(i)+":", flattened)
                }
            }
        } else {
            flatten(value.(map[string]interface{}), key+":", flattened)
        }
    }
}

func encode() string {
    for k, v := range lock.cacheData {
        if _, ok := v.(string); ok {
            jsonString = jsonString+"\""+k+"\":\""+v.(string)+"\","
        } else if _, ok := v.(float64); ok {
            jsonString = jsonString+"\""+k+"\":"+strconv.FormatFloat(v.(float64), 'f', -1, 64)+","
        } else if _, ok := v.(bool); ok {
            jsonString = jsonString+"\""+k+"\":"+strconv.FormatBool(v.(bool))+","
        } else {
            fmt.Println("da fuck happened?")
        }
    }

    if len(jsonString) > 0 {
        jsonString = "{"+jsonString[:(len(jsonString)-1)]+"}"
    }
    return jsonString
}

func show(connection net.Conn, key string) {
    switch key {
        case "KEYS": {
            if len(lock.cacheData) == 0 {
                connection.Write([]byte("NO KEYS TO SHOW YO"))
                return
            } else {
                keys := []string{}
                for k := range lock.cacheData {
                    keys = append(keys, k)
                    sort.Strings(keys)
                }
                connection.Write([]byte(strings.Join(keys, "\n")))
            }
        }
        case "COLLECTIONS": {
            connection.Write([]byte("No collections created yet"))
            return
        }
        default: connection.Write([]byte("Invalid request"))
    }
}

func remove(connection net.Conn, key string) {
    _, ok := lock.cacheData[key]         // check if key is valid
    if ok {
        delete(lock.cacheData, key)
        connection.Write([]byte(key+" has been removed"))
    } else {
        connection.Write([]byte("No key: "+key))
    }  
}

func openDisk(filename string) *os.File {
    disk, err := os.OpenFile(filename, os.O_TRUNC|os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }

    return disk
}

func saveLog(dataInput []byte) { 
    disk := openDisk(LOGFILE)
    defer disk.Close()

    disk.Seek(0,END)
    _, err := disk.Write(dataInput)
    if err != nil {
        log.Fatal(err)
    }
}

func clearLog(filename string) {
    daLog, err := os.OpenFile(filename, os.O_TRUNC, 0666) // Opening it in truncate mode clears the log
    if err != nil {
        log.Fatal(err)
    }
    daLog.Close()
}

// func loadDBOnStart(connection net.Conn) {
    
//     buf := make([]byte, 10000)          // use bytes library for this
//     inputEnd, err := connection.Read(buf)
//     if err == io.EOF {
//         return
//     }

//     dataInput := buf[0:inputEnd]
//     message := string(dataInput)

//     msgSplit := strings.Fields(message)
    
//     if len(msgSplit) < 2 { return }
//     DATABASE = msgSplit[1]

//     load(connection, DATABASE)
//     return
// }

func setDB(connection net.Conn, key string) {
    fmt.Println("old db: ", DATABASE)
    DATABASE := key
    fmt.Println("new db: ", DATABASE)
    
    load(connection, DATABASE)
    return
}

func main() {

    listener, err := net.Listen("tcp", PORT)
    if err != nil {
        log.Fatal(err)
        return
    }

    defer listener.Close()

    // c, err := listener.Accept()   
    // loadDBOnStart(c)
    // fmt.Println("loaded DB")
    // c.Close()
    // //listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
            return
        }

        go echoServer(conn) 
    }
}