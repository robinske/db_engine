
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
    // "sort"
)

// MAKE SURE EACH FUNCTION ONLY DOES ONE THING

var cache = struct{
    sync.RWMutex
    Data map[string]interface{}
}{Data: make(map[string]interface{})}

var keyList = make(map[string]map[string]struct{})
var lkey = ""
var jsonString = ""
var DATABASE string

const (
    PORT = ":4127"
    END = 2
    LOGFILE = "outputs/log.txt"
    BUFFER_SIZE = 1024
)

func echoServer(connection net.Conn) {      // this function does too many things. need to separate it

    for {
        buf := make([]byte, BUFFER_SIZE)    // use bytes library for this
                                            // can increase this number or use a streaming data parser
        inputEnd, err := connection.Read(buf)
        if err == io.EOF {
            return
        }

        dataInput := buf[0:inputEnd]
        message := string(dataInput)

        instruction, key, value := parseRequest(message)
        if instruction == "SET" || instruction == "UPDATE" || instruction == "NESTEDUPDATE" || instruction == "REMOVE" || instruction == "UPDATEINT" || instruction == "ADDTO" {
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
    disk := openDisk(DATABASE)

    disk.Seek(0,END)
    _, err := disk.Write([]byte(data))
    if err != nil {
        log.Fatal(err)
    }

    clearLog(LOGFILE)   // can clear the log once written to stable storage
    disk.Close()
    // connection.Write([]byte("DB has been saved. "))

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
        case "DATABASE:>": {
            DATABASE = key   
            load(connection, DATABASE)
            }
        case "GET": get(connection, key)
        case "KSEARCH": KSearch(connection, key, value)
        // case "VSEARCH": 
        //     value = key
        //     VSearch(connection, value)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "UPDATEINT": 
            intValue, _ := strconv.Atoi(value)
            updateInt(connection, key, intValue)
        case "ADDTO": addto(connection, key, value)
        case "NESTEDUPDATE": updateNested(connection, key, value)
        case "LOAD": load(connection, key)
        case "SHOW": show(connection, key)
        case "REMOVE": remove(connection, key)
        case "QUIT": quit(connection)
        case "SAVE": save(connection)
        case "SEARCH": search(connection, key)
        // case "APPLYLOG": applyLog(connection)
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

// func applyLog(connection net.Conn) {

//     r := bufio.NewReader(LOGFILE)
//     for line, _, err := r.ReadString('\n'); err != io.EOF {
//         connection.Write([]byte(line))
//     }
// }



// func openDisk(filename string) (disk *os.File) {
//     var err error
//     disk, err = os.OpenFile(filename, os.O_TRUNC|os.O_RDWR|os.O_APPEND, 0666)
//     if err != nil {
//         fmt.Fprintf(os.Stderr, "%v\n", err)
//         return
//     }

//     return
// }

// func saveLog(dataInput []byte) { 
//     disk := openDisk(LOGFILE)
//     defer disk.Close()

//     disk.Seek(0,END)
//     _, err := disk.Write(dataInput)
//     if err != nil {
//         log.Fatal(err)
//     }
// }

func KSearch(connection net.Conn, key, value string) string {
    // only works for string values
    
    flattened := make(map[string]interface{})
    flatten(cache.Data, lkey, &flattened)
    values := []string{}


    for k := range flattened {      // NO LONGER HASHING, O(N)
        if strings.Contains(k, key) {
            v := flattened[k]            

            if value == v {
                fmtV := fmt.Sprintf("%v", v)
                values = append(values, k+": "+fmtV) // non efficient space // something to consider if time
            }
        }
    }

    if len(values) == 0 {
        connection.Write([]byte("No values found"))
    } else {
        connection.Write([]byte(strings.Join(values, "\n")+"\n<<Values: "+strconv.Itoa(len(values))+">>"))
    }
    return strings.Join(values, "")
}

// func VSearch(connection net.Conn, value string) {
//     // only works for string values
//     var flatlock = struct{
//         sync.RWMutex
//         flattened map[string]interface{}
//     }{flattened: make(map[string]interface{})}
//     indexValues(flatlock.flattened)
//     values := []string{}

//     if len(values) == 0 {
//         connection.Write([]byte("No values found"))
//     } else {
//         connection.Write([]byte(strings.Join(values, "\n")))
//     }
// }

// func indexValues(flattened map[string]interface{}) (values []string) {
    
//     //flatten values
//     flatten(cache.Data, lkey, &flattened)

//     buckets := make(map[string][]string)

    

//     //create bucket keys for each top level key
//     for key := range cache.Data {
//         buckets[key] = []string{}
//     }

//     values = []string{}

//     for k,v := range flattened {
        
//         switch vv := v.(type) {
//             case string: values = append(values, vv)
//             case float64: values = append(values, strconv.FormatFloat(vv, 'f', -1, 64))
//             case bool: values = append(values, strconv.FormatBool(vv))
//             case nil: values = append(values, "nil")
//         }

//         if strings.Contains(k, ":") {
//             key := strings.Split(k, ":")[0]
//             buckets[key] = values
//         } else {
//             buckets[k] = values
//         }
//     }

//     fmt.Println(buckets)
    

//     //put values (from flattened) into their respective buckets based on the top level key
//     //convert values to strings

//     // 
//     // make an array of values that are stored as pointers to their memory location
//     // keep sorted buckets/indexes (for above) of the values for each root level key
//     // any matching values within the bucket - divide and conquer strategy // look before you leap
//     // return a count of the number of instances of values
//     // once you've then iterated through the keys and hit the count you can stop
//     return
// }

func get(connection net.Conn, key string) {

    cache.RLock()
    value, ok := cache.Data[key]         // check if key is valid
    cache.RUnlock()

    if !ok {
        connection.Write([]byte("Key not valid"))
    } else {
        fmtVV := formatOutput(value, "")
        connection.Write([]byte(fmtVV))
    }
}

func formatOutput(vInterface interface{}, fmtValue string) string {
    switch value := vInterface.(type) {
        case string: fmtValue = "\""+value+"\""
        case float64: fmtValue = fmt.Sprintf("%v", value)
        case int: fmtValue = fmt.Sprintf("%v", value)
        case bool: fmtValue = fmt.Sprintf("%v", value)
        case nil: fmtValue = fmt.Sprintf("%v", value)
        case []interface{}:
            tmp := []string{}
            for _, arrayElem := range value {
                tmpfmtValue := formatOutput(arrayElem, fmtValue)
                tmp = append(tmp, tmpfmtValue)
            }
            fmtValue = "["+strings.Join(tmp, ",\n")+"]"
        case map[string]interface{}:
            tmp := []string{}
            for k, v := range value {
                tmpfmtValue := "\""+k+"\""+": "+formatOutput(v, fmtValue)
                tmp = append(tmp, tmpfmtValue)
            }
            fmtValue = "{"+strings.Join(tmp, ", ")+"}"
        default: fmt.Println("Error Occured")
    }
    return fmtValue
}

func search(connection net.Conn, key string) {
    
    flattened := make(map[string]interface{})
    flatten(cache.Data, lkey, &flattened)
    values := []string{}

    for k := range flattened {      // NO LONGER HASHING, O(N)^M, where M is length of key
        if strings.Contains(k, key) {
            v := flattened[k]
            fmtV := fmt.Sprintf("%v", v)
            values = append(values, k+": "+fmtV)
        }
    }
    if len(values) == 0 {
        connection.Write([]byte("No values found"))
    } else {
        connection.Write([]byte(strings.Join(values, "\n")))
            //+"\n<<Values: "+strconv.Itoa(len(values))+">>"))
    }
}

func set(connection net.Conn, key, value string) {

    _, ok := cache.Data[key]         // check if key is valid
    if ok {
        connection.Write([]byte(key+" already added. To modify, UPDATE key"))
    } else {
        cache.Lock()
        cache.Data[key] = value
        cache.Unlock()
        connection.Write([]byte("Added "+key+":"+value))
    }
}

// have a way to update this with nested vales. recognize nested dictionaries with a / and lists with /# or something???

func updateNested(connection net.Conn, key, value string) {

    keyList := strings.Split(key, "/")

    if len(keyList) == 0 { 
        connection.Write([]byte("Update failed"))
        return 
    }
    collection := keyList[0]

    if len(keyList) == 1 { 
        connection.Write([]byte("Update failed"))
        return 
    }
    uniqueKey := keyList[1]

    if len(keyList) == 2 { 
        connection.Write([]byte("Update failed"))
        return 
    }
    keyToUpdate := keyList[2]


    var update string


    for k, vInterface := range cache.Data {
        if k == collection {
            switch v := vInterface.(type) {
            case []interface{}:
                for i, element := range v {
                    fmtElem := fmt.Sprintf("%v", element)

                    // instead of this add unique keys
                    // like wayyy slow yo.

                    if strings.Contains(fmtElem, uniqueKey) {

                        if _, ok := cache.Data[k].([]interface{})[i].(map[string]interface{}); ok {

                            oldValue := cache.Data[k].([]interface{})[i].(map[string]interface{})[keyToUpdate]
                            cache.Data[k].([]interface{})[i].(map[string]interface{})[keyToUpdate] = value // this will give you the list element you need to append
                            update = fmt.Sprintf("%v Index: %v from %v to %v", k+" "+uniqueKey+" "+keyToUpdate, strconv.Itoa(i), oldValue, value)
                            connection.Write([]byte("Updated "+update))
                            return
                        }
                    }
                }
            case map[string]interface{}: update = "hit the map"
            default: update = "nothing"
            }
        }
    }
    if update != "" {
        connection.Write([]byte("Updated "+update))
    } else {
        connection.Write([]byte("Key string not valid. Nothing updated"))
    }
}

func update(connection net.Conn, key, value string) {
    
    _, ok := cache.Data[key]        // check if key is valid
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {
        cache.Lock()
        cache.Data[key] = value     // overwrite
        cache.Unlock()
        connection.Write([]byte("Updated "+key+":"+value))
    }
}

func updateInt(connection net.Conn, key string, value int) {
    
    _, ok := cache.Data[key]        // check if key is valid
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {
        cache.Lock()
        cache.Data[key] = value     // overwrite
        cache.Unlock()
        strValue := strconv.Itoa(value)
        connection.Write([]byte("Updated "+key+":"+strValue))
    }
}

func addto(connection net.Conn, key, value string) {
    
    _, ok := cache.Data[key]        // check if key is valid
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {

        oldValue := cache.Data[key]
        tmpnewValue := formatOutput(oldValue, "")
        lastChar := tmpnewValue[len(tmpnewValue)-1:]
        tmpnewValue = tmpnewValue[0:len(tmpnewValue)-1]
        newValue := tmpnewValue+", "+value+lastChar
        tmpMap := make(map[string]interface{})        
        decoded := decodeJSONArray([]byte(newValue))
        tmpMap[key] = decoded
        cache.Lock()
        insert(tmpMap, cache.Data)
        cache.Unlock()
        connection.Write([]byte("Updated "+key+":"+newValue))
    }
    save(connection)
}

func removefrom(connection net.Conn, key, value string) {
    
    _, ok := cache.Data[key]        // check if key is valid
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {

        oldValue := cache.Data[key]
        tmpnewValue := formatOutput(oldValue, "")
        lastChar := tmpnewValue[len(tmpnewValue)-1:]
        tmpnewValue = tmpnewValue[0:len(tmpnewValue)-1]
        newValue := tmpnewValue+", "+value+lastChar
        tmpMap := make(map[string]interface{})        
        decoded := decodeJSONArray([]byte(newValue))
        tmpMap[key] = decoded
        cache.Lock()
        insert(tmpMap, cache.Data)
        cache.Unlock()
        connection.Write([]byte("Updated "+key+":"+newValue))
    }
    save(connection)
}

func decodeJSONArray(encodedJSON []byte) []interface{} {
    decoded := []interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
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

        if string(fileContents) == "" {
            connection.Write([]byte("No existing data to load."))
        } else {
            decoded := decodeJSON(fileContents)

            cache.Lock()
            insert(decoded, cache.Data)
            cache.Unlock()

            connection.Write([]byte("Loaded "+filename+" to collection "))
        } 
    } 
}


// func cacheKeys(nested map[string]interface{}) {
//     for key, value := range nested {
//         if _, ok := value.([]interface{}); ok {
//             tmpArray := map[string]struct{}{}
//             for _, k := range value.([]interface{}) {
//                 if _, ok := k.(map[string]interface{}); ok {
//                     for kk := range k.(map[string]interface{}) {
//                         tmpArray[kk] = struct{}{}
//                     }
//                 }
//             }
//             keyList[key] = tmpArray    
//         } else {
//             fmt.Println("Nothing did")
//         }
//     }
// }

// func cacheKeys(nested, cacheKeys map[string]interface{}) map[string]struct{} {

//     for key, value := range nested {
//         switch v := value.(type) {
//             // case float64: nested[key] = v
//             // case string: nested[key] = strings.ToUpper(v)
//             // case bool: nested[key] = v
//             // case nil: nested[key] = nil
//             case map[string]interface{}:  
//                 tmp := map[string]struct{}{}
//                 keyList[key] = cacheKeys(v, tmp)               
//             case []interface{}:
//                 tmpArray := []interface{}{}
//                 for _, val := range v {
//                     if _, ok := val.(map[string]interface{}); ok {
//                         tmp := map[string]interface{}{} 
//                         tmpArray = append(tmpArray, insert(val.(map[string]interface{}), tmp))
//                     } else {
//                         tmpArray = append(tmpArray, val)
//                     }
//                 }
//                 nested[key] = tmpArray
//             default:
//                 fmt.Println("Missed something", v)
//         }
//     }
//     return cacheKeys
// }

func decodeJSON(encodedJSON []byte) map[string]interface{} {
    decoded := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

func insert(inputMap, Data map[string]interface{}) map[string]interface{} {

    for k, value := range inputMap {
        key := strings.ToUpper(k)
        switch v := value.(type) {
            case float64: Data[key] = v
            case string: Data[key] = strings.ToUpper(v)
            case bool: Data[key] = v
            case nil: Data[key] = "nil"
            case map[string]interface{}:  
                tmp := map[string]interface{}{}
                Data[key] = insert(v, tmp)
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
                Data[key] = tmpArray
            default:
                fmt.Println("Missed something", v)
        }
    }
    return Data
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
        } else if value == nil {
            (*flattened)[key] = nil
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
    for key, value := range cache.Data {
        data := "\""+key+"\": "+formatOutput(value, "")
        tmp = append(tmp, data)
    }
    
    // WHEN YOU ADDTO IT ADDS THAT AS A STRING AND ADDS EXTRA QUOTES -- DO SOMETHING ABOUT THIS // RE-ADD INTO THE DICTIONARY
    jsonString = "{"+strings.Join(tmp, ", ")+"}"

    return jsonString
}

func show(connection net.Conn, key string) {
    switch key {
        case "COLLECTIONS": {
            if len(cache.Data) == 0 {
                connection.Write([]byte("NO COLLECTIONS AD`D"))
                return
            } else {
                for k := range cache.Data {
                    connection.Write([]byte(k+"\n"))
                }
            }
        }
        case "DATABASE": connection.Write([]byte(DATABASE))
        case "DATA": {
            if len(cache.Data) == 0 {
                connection.Write([]byte("NO DATA ADDED"))
                return
            } else {
                // fmtData := fmt.Sprintf("%v", cache.Data)
                fmtData := formatOutput(cache.Data, "")
                connection.Write([]byte(fmtData))
            }
        }
        default: connection.Write([]byte("Invalid request"))
    }
}

func remove(connection net.Conn, key string) {
    _, ok := cache.Data[key]         // check if key is valid
    if ok {
        delete(cache.Data, key)
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

func openLog(filename string) (disk *os.File) {
    var err error
    disk, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return
    }

    return
}

func saveLog(dataInput []byte) { 
    disk := openLog(LOGFILE)
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