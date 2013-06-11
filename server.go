
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
)

var cache = struct{ // GLOBAL MAP SERVES AS IN MEMORY DATA STORE
    sync.RWMutex    // LOCK PREVENTS DUAL ACCESS
    Data map[string]interface{}
}{Data: make(map[string]interface{})}

var keyList = make(map[string]map[string]struct{})
var lkey = ""
var jsonString = ""
var DATABASE string
var counter = 0

const (
    PORT = ":4127"
    END = 2
    LOGFILE = "outputs/log2.txt"
    BUFFER_SIZE = 1024
)

func isInArray(element string, list []string) bool {
    for _, b := range list {
        if b == element {
            return true
        }
    }

    return false
}

func runLog(instruction string, dataInput []byte) {

    logInstructions := []string{"SET", "UPDATE", "NESTEDUPDATE", "REMOVE", "UPDATEINT","ADDTO"}

    if isInArray(instruction, logInstructions) {
        stringData := "\n"+string(dataInput)
        saveLog([]byte(stringData))
    }
}

func dispatch(connection net.Conn) {

    for {
        buf := make([]byte, BUFFER_SIZE)
        inputEnd, err := connection.Read(buf)
        if err == io.EOF {
            return
        }

        dataInput := buf[0:inputEnd]
        message := string(dataInput)

        instruction, key, value := parseRequest(message)

        runLog(instruction, dataInput)
        
        // if instruction == "SET" || instruction == "UPDATE" || instruction == "NESTEDUPDATE" || instruction == "REMOVE" || instruction == "UPDATEINT" || instruction == "ADDTO" {
        //     stringData := "\n"+string(dataInput)
        //     saveLog([]byte(stringData))
        // }

        callCacheData(connection, instruction, key, value)

        fmt.Printf("Server received: %s\n", message)
    }
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
        case "SEARCHBYKEY": searchByKey(connection, key, value)
        case "SET": set(connection, key, value)
        case "UPDATE":  update(connection, key, value)
        case "UPDATEINT": 
            intValue, _ := strconv.Atoi(value)
            updateInt(connection, key, intValue)
        // case ""
        case "ADDTO": addto(connection, key, value)
        case "NESTEDUPDATE": updateNested(connection, key, value)
        case "LOAD": load(connection, key)
        case "SHOW": show(connection, key)
        case "REMOVE": remove(connection, key)
        case "QUIT": quit(connection)
        case "SEARCH": search(connection, key)
        case "APPLYLOG": connection.Write([]byte("Updated with most recent log"))
        case "SAVE": save()
        case "CLEARLOG": clearLog(LOGFILE)
        default: connection.Write([]byte("Instruction not recognized"))
    }
}

//////////////////////// READ OPERATIONS ////////////////////////////////////

func show(connection net.Conn, key string) {
    switch key {
        case "COLLECTIONS": {
            if len(cache.Data) == 0 {
                connection.Write([]byte("NO COLLECTIONS ADDED"))
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
                fmtData := formatOutput(cache.Data, "")
                connection.Write([]byte(fmtData))
            }
        }
        default: connection.Write([]byte("Invalid request"))
    }
}


func get(connection net.Conn, key string) {

    cache.RLock()
    value, ok := cache.Data[key]
    cache.RUnlock()

    if !ok {
        connection.Write([]byte("Key not valid"))
    } else {
        fmtVV := formatOutput(value, "")
        connection.Write([]byte(fmtVV))
    }
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
        connection.Write([]byte(strings.Join(values, "\n")+"\n<<Values: "+strconv.Itoa(len(values))+">>"))
    }
}

func searchByKey(connection net.Conn, key, value string) string {
    // only works for string values
    
    flattened := make(map[string]interface{})
    flatten(cache.Data, lkey, &flattened)
    values := []string{}


    for k := range flattened {
        if strings.Contains(k, key) {
            v := flattened[k]            

            if value == v {
                fmtV := fmt.Sprintf("%v", v)
                values = append(values, k+": "+fmtV) // non efficient space // something to consider
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

//////////////////////// MUTATING OPERATIONS ////////////////////////////////////

func set(connection net.Conn, key, value string) {

    _, ok := cache.Data[key]
    if ok {
        connection.Write([]byte(key+" already added. To modify, UPDATE key"))
    } else {
        cache.Lock()
        cache.Data[key] = value
        cache.Unlock()
        connection.Write([]byte("Added "+key+":"+value))
    }
}


func update(connection net.Conn, key, value string) {
    
    _, ok := cache.Data[key]
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {
        cache.Lock()
        cache.Data[key] = value
        cache.Unlock()
        connection.Write([]byte("Updated "+key+":"+value))
    }
}

func updateInt(connection net.Conn, key string, value int) {
    
    _, ok := cache.Data[key]
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {
        cache.Lock()
        cache.Data[key] = value
        cache.Unlock()
        strValue := strconv.Itoa(value)
        connection.Write([]byte("Updated "+key+":"+strValue))
    }
}



func updateNested(connection net.Conn, key, value string) {

    keyList := strings.Split(key, "/")

    // NEED AT LEAST 3 VARIABLES TO UPDATED A NESTED VALUE

    if len(keyList) < 3 { 
        connection.Write([]byte("Update failed"))
        return 
    }

    collection := keyList[0]
    uniqueKey := keyList[1]
    keyToUpdate := keyList[2]

    var update string

    for k, vInterface := range cache.Data {
        if k == collection {
            switch v := vInterface.(type) {
            case []interface{}:
                for i, element := range v {
                    fmtElem := fmt.Sprintf("%v", element)
                    if strings.Contains(fmtElem, uniqueKey) {   // could be much more efficient
                        if _, ok := cache.Data[k].([]interface{})[i].(map[string]interface{}); ok {
                            oldValue := cache.Data[k].([]interface{})[i].(map[string]interface{})[keyToUpdate]
                            cache.Data[k].([]interface{})[i].(map[string]interface{})[keyToUpdate] = value // this will give you the list element you need to append
                            update = fmt.Sprintf("%v Index: %v from %v to %v", k+" "+uniqueKey+" "+keyToUpdate, strconv.Itoa(i), oldValue, value)
                            connection.Write([]byte("Updated "+update))
                            return
                        }
                    }
                }
            default: update = "nothing"
            }
        }
    }
    if update != "" {
        connection.Write([]byte("Updated "+update))
    } else {
        connection.Write([]byte("Key string invalid. Nothing updated"))
    }
}

func addto(connection net.Conn, key, value string) {
    
    _, ok := cache.Data[key]
    if !ok {
        connection.Write([]byte(key+" not yet added. Please set"))
    } else {
        // remember what it used to be
        oldValue := cache.Data[key]
        tmpnewValue := formatOutput(oldValue, "")

        // saves info about whether it was an array or object
        lastChar := tmpnewValue[len(tmpnewValue)-1:]
        tmpnewValue = tmpnewValue[0:len(tmpnewValue)-1]

        // append the new value
        newValue := tmpnewValue+", "+value+lastChar
        
        // reinsert into cache as nested objects
        tmpMap := make(map[string]interface{})        
        decoded := decodeJSONArray([]byte(newValue))
        tmpMap[key] = decoded
        cache.Lock()
        insert(tmpMap, cache.Data)
        cache.Unlock()
        
        connection.Write([]byte("Updated "+key+":"+newValue))
    }
    // save(connection) // necessary?
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


//////////////////////// CACHE INPUT INTERACTIONS ////////////////////////////////////


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
                fmt.Println("Error Occured", v)
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

//////////////////////// FORMATTERS ////////////////////////////////////////

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

func decodeJSON(encodedJSON []byte) map[string]interface{} {
    decoded := map[string]interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

func decodeJSONArray(encodedJSON []byte) []interface{} {
    decoded := []interface{} {}
    err := json.Unmarshal(encodedJSON, &decoded)
    if err != nil {
        log.Fatal(err)
    }
    return decoded
}

//////////////////////// DISK OPERATIONS ////////////////////////////////////////

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

func applyLog(connection net.Conn) {

    buf := make([]byte, BUFFER_SIZE)

    fileContents, err := ioutil.ReadFile(LOGFILE)
    fileString := string(fileContents)

    if err != nil {
        log.Fatal(err)
    }

    fileArray := strings.Split(fileString, "\n")
    fileArray = append(fileArray, "CLEARLOG")

    for _, line := range fileArray {
        if line != "" {
            connection.Write([]byte(line))
            inputEnd, err := connection.Read(buf[:])
            if err != nil {
                return
            }
            fmt.Printf("%s\n", string(buf[0:inputEnd]))
        }
    }

    connection.Write([]byte("Updated with most recent log"))

}

func saveLog(dataInput []byte) { 

    counter++

    disk := openLog(LOGFILE)
    defer disk.Close()

    disk.Seek(0,END)

    _, err := disk.Write(dataInput)
    if err != nil {
        log.Fatal(err)
    }

    if counter%5 == 0 {
        save()
    }
}

func clearLog(filename string) {

    daLog, err := os.OpenFile(filename, os.O_TRUNC, 0666) // Opening it in truncate mode clears the log
    
    if err != nil {
        log.Fatal(err)
    }

    daLog.Close()
}

func quit(connection net.Conn) {

    if DATABASE != "" {
        save()
        connection.Write([]byte("Program exiting"))       
    } else {
        connection.Write([]byte("No database set, changes have not been saved. Program exiting"))
    }

    os.Exit(0)
}

func save() {

    data := encode()
    disk := openDisk(DATABASE)

    disk.Seek(0,END)
    _, err := disk.Write([]byte(data))
    if err != nil {
        log.Fatal(err)
    }

    clearLog(LOGFILE)   // can clear the log once written to stable storage
    disk.Close()

}

func encode() string {
    tmp := []string{}
    for key, value := range cache.Data {
        data := "\""+key+"\": "+formatOutput(value, "")
        tmp = append(tmp, data)
    }
    
    jsonString = "{"+strings.Join(tmp, ", ")+"}"

    return jsonString
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

        go dispatch(conn)     // goroutine manages concurrent processes
    }
}