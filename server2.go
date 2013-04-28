
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



var lkey = ""
var jsonString = ""
var collectionList []string
var DATABASE string

const (
    PORT = ":4127"
    END = 2
    LOGFILE = "outputs/log.txt"
    BUFFER_SIZE = 1e8
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
        case "GETWHERE": getWhere(connection, key, value)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "LOAD": load(connection, key)
        case "SHOW": show(connection, key)
        case "REMOVE": remove(connection, key)
        case "QUIT": quit(connection)
        case "SAVE": save(connection)
        case "RECON": reconstructJSON(connection, lock.nested)
        case "JSON": 
            connection.Write([]byte("see server"))
            fmt.Println(jsonString)
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

func getWhere(connection net.Conn, key, value string) {
    // only works for string values

    values := []string{}

    for k := range lock.nested {      // NO LONGER HASHING, O(N)
        if strings.Contains(k, key) {
            lock.RLock()
            v := lock.nested[k]            
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

// search keys -- need an index of keys - exisitng string function minus the data

func get(connection net.Conn, key string) {

    lock.RLock()
    value, ok := lock.nested[key]         // check if key is valid
    lock.RUnlock()

    if !ok {
        connection.Write([]byte("Key not valid"))
    } else {
        fmtValue := fmt.Sprintf("%v", value)
        fmt.Println(fmtValue)
        connection.Write([]byte(fmtValue))
    }
}

func set(connection net.Conn, key, value string) {

    // make clear for which collection

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

func update(connection net.Conn, key, value string) {
    
    _, ok := lock.nested[key]        // check if key is valid
    if ok {
        lock.Lock()
        lock.nested[key] = value     // overwrite
        lock.Unlock()
        connection.Write([]byte("Updated "+key+":"+value))
    } else {
        connection.Write([]byte(key+" not yet added. Please set"))
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

        insert(decoded, lock.nested)

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
   
    //define a root node / variable - initialize this here
    // root := inputMap

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


                // for _, val := range v {
                //     if _, ok := val.(map[string]interface{}); ok {
                //         nested[key] = insert(val.(map[string]interface{}), root)
                //     } else if _, ok := val.(int); ok {
                //         v = append(v, strings.ToUpper(strconv.Itoa(val.(int))))
                //         lock.nested[key] = v
                //     } else if _, ok := val.(string); ok {
                //         v = append(v, strings.ToUpper(val.(string)))
                //         lock.nested[key] = v
                //     } else {
                //         fmt.Println("yikes")
                //     }
                // }

// if value of key a basic type (float, string, bool or nil, insert key : value into the dictionary
// if value of key is a map, RECURSE: insert another map as value of key
// if value is an array, RECURSE OVER ELEMENTS, insert ARRAY as value of key




// func insert(inputJSON map[string]interface{}, nested *map[string]interface{}) {



//     for k, value := range inputJSON {
//         key := strings.ToUpper(k)
//         switch v := value.(type) {
//             case float64: (*nested)[key] = v
//             case string: (*nested)[key] = strings.ToUpper(v)
//             case bool: (*nested)[key] = v
//             case map[string]interface{}: 
//                 insert(v, nested)
//                 (*nested)[key] = v
//             case []interface{}: 
//                 for _, val := range v {
//                     if _, ok := val.(map[string]interface{}); ok {
//                         //(*nested)[key] = val
//                         insert(val.(map[string]interface{}), nested)
//                     } else if _, ok := val.(int); ok {
//                         v = append(v, strings.ToUpper(strconv.Itoa(val.(int))))
//                         (*nested)[key] = v
//                     } else if _, ok := val.(string); ok {
//                         v = append(v, strings.ToUpper(val.(string)))
//                         (*nested)[key] = v
//                     } else {
//                         fmt.Println("yikes")
//                     }
//                 }
//             default:
//                 fmt.Println("ABLAKJFALSKDFJ WHAT DID I MISS", v)
//         }
//     }
// }

// if value of key a basic type (float, string, bool or nil, insert key : value into the dictionary
// if value of key is a map, RECURSE: insert another map as value of key
// if value is an array, RECURSE OVER ELEMENTS, insert ARRAY as value of key

func reconstructJSON(connection net.Conn, inputMap map[string]interface{}) string {

    // fmtList := []string{}

    // for key, value := range inputMap {
    //     fmtJSON := fmt.Sprintf("%v: \n\t%v\n", key, value)
    //     // fmtList = append(fmtList, fmtSomething)
    //     connection.Write([]byte(fmtJSON))
    // }

    for key, value := range inputMap {
        switch v := value.(type) {
            case string: jsonString = jsonString+"\""+key+"\":\""+value.(string)+"\",\n"
            case float64: jsonString = jsonString+"\""+key+"\":"+strconv.FormatFloat(value.(float64), 'f', -1, 64)+",\n"
            case bool: jsonString = jsonString+"\""+key+"\":"+strconv.FormatBool(value.(bool))+",\n"
            case []interface{}:
                connection.Write([]byte("list triggered"+key))
                jsonString = jsonString+"\""+key+"\":"+"\n["
                for _, elem := range v {
                    if _, ok := elem.(map[string]interface{}); ok {
                        jsonString = jsonString+"\t{\n"
                        reconstructJSON(connection, elem.(map[string]interface{}))
                        jsonString = jsonString+"},\n"
                    } else {
                        connection.Write([]byte("hmmmmm"))
                    }
                }
                jsonString = jsonString+"]\n"
            case map[string]interface{}:
                jsonString = jsonString+"\t"+"\""+key+"\":\n\t{\n"
                reconstructJSON(connection, v)

            default: 
                fmtIt := fmt.Sprintf("%T", key)
                connection.Write([]byte("WHAT DID I MISS"+key+fmtIt))
        }
    }

    // fmt.Println(jsonString)
    // connection.Write([]byte(jsonString))

    // for _, elem := range fmtList {
    //     fmt.Println(elem, "\n")
    // }
    
    //connection.Write([]byte("see server"))

    return jsonString

}

func encode() string {
    for k, v := range lock.nested {
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
                fmtData := fmt.Sprintf("%v", lock.nested)
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