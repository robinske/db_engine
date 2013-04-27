
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
    nested map[string]interface{}
}{nested: make(map[string]interface{})}

// other locking options have a single goroutine that's responsible for applying changes ot the database only one channel to send the mutations -- will manage a queue of requests

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
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

func getWhere(connection net.Conn, key, value string) {
   
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


// func typeSwitch(connection net.Conn, value map[string]interface{}) {
    
//     // for
//     //     switch val := value.(type) {
//     //         case string: data = append(data, (v))
//     //         case float64: data = append(data, (strconv.FormatFloat(v, 'f', -1, 64)))
//     //         case bool: data = append(data, (strconv.FormatBool(v)))
            
//     //                 typeSwitch(connection, v.(map[string]interface{}))
//     //             }
//     //     }

//     for k, v := range value {
//         switch vv := v.(type) {
//         case string:
//             fmt.Println(k, "is string", vv)
//         case int:
//             fmt.Println(k, "is int", vv)
//         case []interface{}:
//             fmt.Println(k, "is an array:")
//             for i, u := range vv {
//                 fmt.Println(i, u)
//             }
//         default:
//             fmt.Println(k, "is of a type I don't know how to handle")
//         }
//     }

//                 // case []interface{}: 
//                 //     for _, val := range v {
//                 //         if _, ok := val.(string); ok {
//                 //             v = append(v, strings.ToUpper(val.(string)))
//                 //             (*nested)[key] = v
//                 //         } else if _, ok := val.(int); ok {
//                 //             v = append(v, strings.ToUpper(strconv.Itoa(val.(int))))
//                 //             (*nested)[key] = v
//                 //         } else {
//                 //             insert(val.(map[string]interface{}), nested)
//                 //             (*nested)[key] = v
//                 //         }
//                 //     }
//                 // case map[string]interface{}: 
//                 //     (*nested)[key] = v
//                 //     insert(v, nested)

//     sort.Strings(data)
//     connection.Write([]byte(strings.Join(data, "\n")))
// }

func get(connection net.Conn, key string) {

    lock.RLock()
    value, ok := lock.nested[key]         // check if key is valid
    lock.RUnlock()

    if !ok {
        connection.Write([]byte("Key not valid"))
    } else {
        //connection.Write()
        switch v := value.(type) {
            case string: connection.Write([]byte(v))
            case bool: connection.Write([]byte(strconv.FormatBool(v)))
            case float64: connection.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
            case nil: connection.Write([]byte("nil"))
            

            case []interface{}:
                for _, vv := range v {
                    typeSwitch(connection, key, vv)
                }
            case map[string]interface{}:
                // for k, vv := range v {

                // }



            default: {
                //typeSwitch(connection, key, value)
                // if len(newdata) != 0 {
                //     sort.Strings(newdata)
                //     connection.Write([]byte(strings.Join(newdata, "\n")))
                // } else {
                connection.Write([]byte("didn't work"))
                // }
            }
        }
    }
}

func typeSwitch(connection net.Conn, key string, value interface{}) []string {
    
    newdata := []string{}
    switch v := value.(type) {
        case string: newdata = append(newdata, v)
        case bool: newdata = append(newdata, strconv.FormatBool(v))
        case float64: newdata = append(newdata, strconv.FormatFloat(v, 'f', -1, 64))
        case nil: newdata = append(newdata, "nil")
        // case []interface{}:
        //     fmt.Println("got here tooooo")
        //     for _, elem := range v {
        //         fmt.Println(elem)
        //         tempdata := typeSwitch(connection, key, elem)
        //         newdata = append(newdata, tempdata)
        //         return newdata                    
        //     }
        //     return newdata
        // case map[string]interface{}:
        //     tempdata := typeSwitch(connection, key, v)
        //     newdata = append(newdata, tempdata)
        //     return newdata
        default: 
            // RECONSTRUCT JSON - SEND OVER THIS JSON STRING
            connection.Write([]byte("Invalid JSON fail"))
            // return
    }
    return newdata
}
    
        // check for end value types: string, bool, float64, and nil
        // check for []interface{} - recurse on list elements
        // check for map[string]interface{} - recurse on value

        




    //     if _, ok := value.(string); ok {
    //         connection.Write([]byte(value.(string)))
    //     } else if _, ok := value.(float64); ok {
    //         connection.Write([]byte(strconv.FormatFloat(value.(float64), 'f', -1, 64)))
    //     } else if _, ok := value.([]string); ok {
    //         connection.Write([]byte(strings.Join(value.([]string)[:], " ")))
    //     } else if _, ok := value.(map[string]interface{}); ok {
    //         // typeSwitch(connection, value.(map[string]interface{}))
    //         connection.Write([]byte("see server"))
    //         for k, v := range value.(map[string]interface{}) {
    //             fmt.Println(k, ": ", v)
    //         }
    //         // fmt.Println(value)
    //     } else {
    //         // data := []string{}
    //         for k, v := range value.(map[string]interface{}) {
    //             fmt.Println(k, ": ", v)
    //         }
    //         connection.Write([]byte("see server"))
    //         fmt.Println(value)
    //     }
    // } else {
    //     connection.Write([]byte("No values found"))
    // }
// }

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

func load(connection net.Conn, filename string) { //, optCollection...string
    
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
        
        mappedJSON := decodeJSON(fileContents)

        insert(mappedJSON, &lock.nested)

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

func insert(inputJSON map[string]interface{}, nested *map[string]interface{}) {
    for k, value := range inputJSON {
        key := strings.ToUpper(k)
        switch v := value.(type) {
            case float64: (*nested)[key] = v
            case string: (*nested)[key] = strings.ToUpper(v)
            case bool: (*nested)[key] = v
            case []interface{}: 
                for _, val := range v {
                    if _, ok := val.(string); ok {
                        v = append(v, strings.ToUpper(val.(string)))
                        (*nested)[key] = v
                    } else if _, ok := val.(int); ok {
                        v = append(v, strings.ToUpper(strconv.Itoa(val.(int))))
                        (*nested)[key] = v
                    } else {
                        insert(val.(map[string]interface{}), nested)
                        (*nested)[key] = v
                    }
                }
            case map[string]interface{}: 
                (*nested)[key] = v
                insert(v, nested)
            default:
                fmt.Println("ABLAKJFALSKDFJ WHAT DID I MISS", v)
        }
    }
}

// func reconstructJSON(nested interface{}) string {

// }

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
                data := []string{}
                for key, value := range lock.nested {
                    switch v := value.(type) {
                        case string: data = append(data, (key+": "+v))
                        case float64: data = append(data, (key+": "+strconv.FormatFloat(v, 'f', -1, 64)))
                        case bool: data = append(data, (key+": "+strconv.FormatBool(v)))
                        default: fmt.Println("huhhhh???")                
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