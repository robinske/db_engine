
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

var lock = struct{
    sync.RWMutex
    nested map[string]interface{}
}{nested: make(map[string]interface{})}

var flatlock = struct{
    sync.RWMutex
    flattened map[string]interface{}
}{flattened: make(map[string]interface{})}

var lkey = ""
var jsonString = ""
var collectionList []string
var DATABASE string

const (
    PORT = ":4127"
    END = 2
    LOGFILE = "outputs/log.txt"
    BUFFER_SIZE = 1e9
)

func echoServer(connection net.Conn) {      // this function does too many things. need to separate it

    for {
        buf := make([]byte, BUFFER_SIZE)          // use bytes library for this
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

        callnested(connection, instruction, key, value)

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

func callnested(connection net.Conn, instruction, key string, optionalValue...string) {

    value := strings.Join(optionalValue[:], " ")

    switch instruction {
        case "DATABASE:>": {
            DATABASE = key   
            load(connection, DATABASE)
            }
        case "GET": get(connection, key)
        case "KSEARCH": KSearch(connection, key, value)
        case "VSEARCH": 
            value = key
            VSearch(connection, value)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "LOAD": load(connection, key)
        case "SHOW": show(connection, key)
        case "REMOVE": remove(connection, key)
        case "QUIT": quit(connection)
        case "SAVE": save(connection)
        case "SEARCH": search(connection, key)
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

func KSearch(connection net.Conn, key, value string) {
    // only works for string values
    flatten(lock.nested, lkey, &flatlock.flattened)
    values := []string{}


    for k := range flatlock.flattened {      // NO LONGER HASHING, O(N)
        if strings.Contains(k, key) {
            // lock.RLock()
            v := flatlock.flattened[k]            
            // lock.RUnlock()

            if value == v {
                values = append(values, k+": "+v.(string)) // non efficient space complexity // something to consider if time
            }
        }
    }

    if len(values) == 0 {
        connection.Write([]byte("No values found"))
    } else {
        connection.Write([]byte(strings.Join(values, "\n")+"\n<<Values: "+strconv.Itoa(len(values))+">>"))
    }
}

func VSearch(connection net.Conn, value string) {
    // only works for string values
    indexValues(flatlock.flattened)
    values := []string{}

    if len(values) == 0 {
        connection.Write([]byte("No values found"))
    } else {
        connection.Write([]byte(strings.Join(values, "\n")))
    }
}

func indexValues(flattened map[string]interface{}) {
    
    flatten(lock.nested, lkey, &flattened)

    buckets := make(map[string][]string)

    for key := range lock.nested {
        buckets[key] = []string{}
    }



    //flatten values
    //create bucket keys for each top level key
    //
    //put values (from flattened) into their respective buckets based on the top level key
    //convert values to strings
    //

    // 
    // make an array of values that are stored as pointers to their memory location
    // keep sorted buckets/indexes (for above) of the values for each root level key
    // any matching values within the bucket - divide and conquer strategy // look before you leap
    // return a count of the number of instances of values
    // once you've then iterated through the keys and hit the count you can stop
}

func get(connection net.Conn, key string) {

    lock.RLock()
    value, ok := lock.nested[key]         // check if key is valid
    lock.RUnlock()

    if !ok {
        connection.Write([]byte("Key not valid"))
    } else {
        fmtVV := formatOutput(value, "")
        connection.Write([]byte(fmtVV))
    }
}

func formatOutput(value interface{}, fmtVV string) string {
    switch vv := value.(type) {
        case string: fmtVV = "\""+vv+"\""
        case float64: fmtVV = fmt.Sprintf("%v", vv)
        case bool: fmtVV = fmt.Sprintf("%v", vv)
        case nil: fmtVV = fmt.Sprintf("%v", vv)
        case []interface{}:
            tmp := []string{}
            for _, val := range vv {
                fmtVal := formatOutput(val, fmtVV)
                tmp = append(tmp, fmtVal)
            }
            fmtVV = "["+strings.Join(tmp, ",\n")+"]"
        case map[string]interface{}:
            tmp := []string{}
            for key, val := range vv {
                fmtVal := "{"+"\""+key+"\""+": "+formatOutput(val, fmtVV)+"}"
                tmp = append(tmp, fmtVal)
            }
            fmtVV = strings.Join(tmp, ", ")
        default: fmt.Println("Error Occured")
    }
    return fmtVV
}

func search(connection net.Conn, key string) {
// show a count???
    flatlock.Lock()
    flatten(lock.nested, lkey, &flatlock.flattened)
    flatlock.Unlock()
    values := []string{}

    for k := range flatlock.flattened {      // NO LONGER HASHING, O(N)^M, where M is length of key
        if strings.Contains(k, key) {
            flatlock.RLock()
            v := flatlock.flattened[k]
            flatlock.RUnlock()
            fmtV := fmt.Sprintf("%v", v)
            values = append(values, k+": "+fmtV)
        }
    }
    if len(values) == 0 {
        connection.Write([]byte("No values found"))
    } else {
        connection.Write([]byte(strings.Join(values, " \n")))
    }
}

func set(connection net.Conn, key, value string) {

    _, ok := lock.nested[key]         // check if key is valid
    if ok {
        connection.Write([]byte(key+" already added. To modify, UPDATE key"))
    } else {
        lock.Lock()
        lock.nested[key] = value
        lock.Unlock()
        connection.Write([]byte("Added "+key+":"+value))
    }
}

// have a way to update this with nested vales. recognize nested dictionaries with a / and lists with /# or something???

func update(connection net.Conn, key, value string) {
    
    _, ok := lock.nested[key]        // check if key is valid
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {
        lock.Lock()
        lock.nested[key] = value     // overwrite
        lock.Unlock()
        connection.Write([]byte("Updated "+key+":"+value))
    }
}

func load(connection net.Conn, filename string) {
    
    if filename == "" {
        connection.Write([]byte("Please enter the data you would like to load"))
        return

    } else {
        
        fileContents, err := ioutil.ReadFile(filename)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%v\n", err)
            connection.Write([]byte("Invalid file"))
            return
        }
        
        decoded := decodeJSON(fileContents)

        lock.Lock()
        insert(decoded, lock.nested)
        lock.Unlock()

        connection.Write([]byte("Loaded "+filename+" to collection "))
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

func insert(inputMap, nested map[string]interface{}) map[string]interface{} {

    for k, value := range inputMap {
        key := strings.ToUpper(k)
        switch v := value.(type) {
            case float64: nested[key] = v
            case string: nested[key] = strings.ToUpper(v)
            case bool: nested[key] = v
            case map[string]interface{}:  
                tmp := map[string]interface{}{}
                nested[key] = insert(v, tmp)
            case []interface{}:
                tmpArray := []interface{}{}
                for _, val := range v {
                    if _, ok := val.(map[string]interface{}); ok {
                        tmp := map[string]interface{}{} 
                        tmpArray = append(tmpArray, insert(val.(map[string]interface{}), tmp))
                    } else {
                        tmpArray = append(tmpArray, val)
                    }
                }
                nested[key] = tmpArray
            default:
                fmt.Println("Missed something", v)
        }
    }
    return nested
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
        } else if _, ok := value.(map[string]interface{}); ok {
            flatten(value.(map[string]interface{}), key+":", flattened)
        } else {
            fmt.Println("failed somehow", key)
        }
    }
}

func encode() string {
    tmp := []string{}
    for key, value := range lock.nested {
        data := "\""+key+"\": "+formatOutput(value, "")
        tmp = append(tmp, data)
    }
    
    jsonString = "{"+strings.Join(tmp, ", ")+"}"

    return jsonString
}

func show(connection net.Conn, key string) {
    switch key {
        case "KEYS": {
            if len(lock.nested) == 0 {
                connection.Write([]byte("NO KEYS TO SHOW YO"))
                return
            } else {
                keys := []string{}
                for k := range lock.nested {
                    keys = append(keys, k)
                    sort.Strings(keys)
                }
                connection.Write([]byte(strings.Join(keys, "\n")))
            }
        }

        case "DATABASE": connection.Write([]byte(DATABASE))
        case "DATA": {
            if len(lock.nested) == 0 {
                connection.Write([]byte("NO DATA TO SHOW YO"))
                return
            } else {
                // fmtData := fmt.Sprintf("%v", lock.nested)
                fmtData := formatOutput(lock.nested, "")
                connection.Write([]byte(fmtData))
            }
        }
        default: connection.Write([]byte("Invalid request"))
    }
}

func remove(connection net.Conn, key string) {
    _, ok := lock.nested[key]         // check if key is valid
    if ok {
        delete(lock.nested, key)
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