
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
type JSON map[string]interface{}

type Tuple struct {
    key string
    value interface{}
}

var lock = struct{
    sync.RWMutex
    cacheData map[string]interface{}
}{cacheData: make(map[string]interface{})}

// other locking options have a single goroutine that's responsible for applying changes ot the database only one channel to send the mutations -- will manage a queue of requests

var flattened = make(map[string]interface{})
var lkey = ""
var jsonString = ""
var collectionList []string
var DATABASE string

const (
    PORT = ":4127"
    END = 2
    LOGFILE = "outputs/log.txt"
)

func echoServer(connection net.Conn) {      // this function does too many things. need to separate it

    for {
        buf := make([]byte, 10000)          // use bytes library for this
                                            // can increase this number or use a streaming data parser
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

        fmt.Printf("Server received: %s\n", message)
    }
}

func quit(connection net.Conn) {

    if DATABASE != "" {
        save(connection)
        connection.Write([]byte("Program exiting"))       
    } else {
        connection.Write([]byte("No database set, changes have not been saved. Program exiting"))
    }

    os.Exit(0)
}

func save(connection net.Conn) {

    data := encode()
    disk := openDisk("outputs/output")

    disk.Seek(0,END)
    _, err := disk.Write([]byte(data))
    if err != nil {
        log.Fatal(err)
    }

    clearLog(LOGFILE)   // can clear the log once written to stable storage
    disk.Close()
    connection.Write([]byte("DB has been saved. "))

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

func callCacheData(connection net.Conn, instruction, key string, optionalValue...string) {

    value := strings.Join(optionalValue[:], " ")

    switch instruction {
        // case "CREATE": {
        //     collection := key
        //     create(connection, collection)
        // }
        case "DATABASE:>": {
            DATABASE = key   
            load(connection, DATABASE)
            }
        case "GET": get(connection, key)
        case "GETWHERE": getWhere(connection, key, value)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "LOAD": load(connection, key)
        case "SHOW": show(connection, key)
        case "REMOVE": remove(connection, key)
        case "QUIT": quit(connection)
        case "SAVE": save(connection)
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

// func create(connection net.Conn, collection string) {
    
//     collectionList = append(collectionList, collection)

//     lock.Lock()    
//     lock.cacheData["<COLLECTIONS>"] = collectionList
//     lock.Unlock()

//     connection.Write([]byte("Created "+collection))
// }

func getWhere(connection net.Conn, key, value string) {
   
    values := []string{}

    for k := range lock.cacheData {      // NO LONGER HASHING, O(N)
        if strings.Contains(k, key) {
            lock.RLock()
            v := lock.cacheData[k]
            lock.RUnlock()

            if value == v {
                values = append(values, k+": "+v.(string))
            }
        }
    }

    if len(values) == 0 {
        connection.Write([]byte("No values found"))
    } else {
        connection.Write([]byte(strings.Join(values, "\n")))
    }
}

// make an array of values that are stored as pointers to their memory location
// keep sorted buckets/indexes (for above) of the values for each root level key
// any matching values within the bucket - divide and conquer strategy // look before you leap
// return a count of the number of instances of values
// once you've then iterated through the keys and hit the count you can stop

func get(connection net.Conn, key string) {

    value, ok := lock.cacheData[key]         // check if key is valid
    values := []string{}

    if ok {
        if _, ok := value.(string); ok {
            connection.Write([]byte(value.(string)))
        } else if _, ok := value.([]string); ok {
            connection.Write([]byte(strings.Join(value.([]string)[:], " ")))
        }
    } else {
        for k := range lock.cacheData {      // NO LONGER HASHING, O(N)^M, where M is length of key
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

    // make clear for which collection

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
    
    _, ok := lock.cacheData[key]        // check if key is valid
    if ok {
        lock.Lock()
        lock.cacheData[key] = value     // overwrite
        lock.Unlock()
        connection.Write([]byte("Updated "+key+":"+value))
    } else {
        connection.Write([]byte(key+" not yet added. Please set"))
    }
}

func load(connection net.Conn, filename string) { //, optCollection...string
    
    if filename == "" {
        connection.Write([]byte("Please enter the data you would like to load"))
        return

    // } else if len(optCollection) == 0 {
    //     connection.Write([]byte("No collection entered, please specify collection"))
    //     return

    } else {
        
        // collection := strings.Join(optCollection[:], " ")
        // create(connection, collection)

        fileContents, err := ioutil.ReadFile(filename)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%v\n", err)
            connection.Write([]byte("Invalid file"))
            return
        }
        
        mappedJSON := decodeJSON(fileContents)

        flatten(mappedJSON, lkey, &flattened)
        
        for k,v := range flattened {
            k = strings.ToUpper(k)
            // k = "<"+strings.ToUpper(collection)+">"+k
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
        connection.Write([]byte("Loaded "+filename+" to collection ")) // +collection))
    } 
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
        } else if _, ok := v.([]string); ok {
            valuestring := ""
            for _, value := range v.([]string) {
                valuestring = valuestring+"\""+value+"\","
                fmt.Println("adding value", value)
            }
            valuestring = valuestring[:(len(valuestring)-1)]
            jsonString = jsonString+"\""+k+"\":"+"["+valuestring+"],"
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
        // case "COLLECTIONS": {
        //     connection.Write([]byte("No collections created yet"))
        //     return
        // }
        case "DATABASE": connection.Write([]byte(DATABASE))
        case "DATA": {
            if len(lock.cacheData) == 0 {
                connection.Write([]byte("NO DATA TO SHOW YO"))
                return
            } else {
                data := []string{}
                for key, value := range lock.cacheData {
                    if _, ok := value.(string); ok {
                        data = append(data, (key+": "+value.(string)))     
                    } else if _, ok := value.(float64); ok {
                        data = append(data, (key+": "+strconv.FormatFloat(value.(float64), 'f', -1, 64)))
                    } else if _, ok := value.(bool); ok {
                        data = append(data, (key+": "+strconv.FormatBool(value.(bool))))
                    } else {
                        fmt.Println("huhhhh???")
                        return
                    }
                }
                sort.Strings(data)
                connection.Write([]byte(strings.Join(data, "\n")))
            }
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

func openDisk(filename string) (disk *os.File) {
    var err error
    disk, err = os.OpenFile(filename, os.O_TRUNC|os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return
    }

    return
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

func main() {

    listener, err := net.Listen("tcp", PORT)
    if err != nil {
        log.Fatal(err)
        return
    }

    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
            return
        }

        go echoServer(conn) 
    }
}