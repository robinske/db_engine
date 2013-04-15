// CLIENT DIALS / WRITES

package main

  import (
    "net"
    "io"
    "log"
    "fmt"
    "strings"
    "bufio"
    "os"
  )

type cacheData map[string]string
// type dataObject

func reader(r io.Reader) {
  buf := make([]byte, 1024) // makes a list of bytes / why 1024? why a list? whyyyyy
  for {
    n, err := r.Read(buf[:]) // calls the read method on r (which is passed in but what is r?? sets this instance to n
    if err != nil {
      return
    }
    println("Client dialed:", string(buf[0:n]))
  }
}

func talkToDictionary(instruct, key, value string, dictionary cacheData) (cacheData) {
  
  
  // ints := []int{1,2,3}
  // values := []string{instruct, key, value}

  // dict2 := make(map[ints int] values string)
  // fmt.Println(dict2)

  switch instruct {
    case "GET": get(key, instruct, dictionary)
    case "PUT": put(key, value, dictionary)
    //case "SAVE": save(key, value, instruct, dictionary)
    default: fmt.Println("try again idiot")
  }

  fmt.Println("YOU'VE TALKED TO THE DB")
  fmt.Println(dictionary)
  return dictionary
}

func get(key, instruct string, dictionary cacheData) {
  getKey := key
  getInstruct := instruct

  fmt.Printf("You called the %s function\n", getInstruct)
  fmt.Printf("We'll get the value of key %s\n", getKey)
  fmt.Printf(dictionary[getKey]) ////// WHYYYYYY
  fmt.Println("should have returned the value")
}

func put(key, value string, dictionary cacheData) {
  // fmt.Printf("You called the %s function\n", instruct)
  putKey := key
  putValue := value

  fmt.Printf("We'll put %s:%s in the dictionary\n", putKey, putValue)
  dictionary[putKey] = putValue
  //fmt.Println(data[key]) // should print value
}

//func save(key, value, instruct string) {
//   fmt.Printf("You called the %s function\n", instruct)
//   fmt.Printf("We'll save %s:%s from dictionary to disk\n", key, value)
// }

func parseRequest(message string) (string, cacheData) {
  msgSplit := strings.Split(message, " ")

  //fmt.Printf("%s", message)
  //fmt.Printf("%s", msgSplit)

  instruct := msgSplit[0]
  key := msgSplit[1]
  value := msgSplit[2]
  // value := msgSplit[2:]
  // if len(value) == 0 {
  //   return
  // }

  // if len(msgSplit) < 2 {
    // handle error
  // }
  // database := "database" // could be a "FROM" statement later
  var dictionary = make(cacheData) //{key:value}
  fmt.Println("printing base dictionary", dictionary)

  talkToDictionary(instruct, key, value, dictionary)

  return value, dictionary
}

func main() {
  c, err := net.Dial("tcp", ":4127") // sets a connection, c, to the port 4127
  if err != nil {
    log.Fatal(err)
  }
  defer c.Close()

  go reader(c) // concurrent process / Goroutine
  //go talkToDB()
  
  for {
    input := bufio.NewReader(os.Stdin)
    rawMessage, err := input.ReadString('\n') // message resets with each new input
    message := strings.ToUpper(rawMessage) // capitalize/normalize input
    if err != nil {
      log.Fatal(err)
    }
    if message != "" {
        parseRequest(message)
        // DO THE SHIT HERE
        _,err := c.Write([]byte(message))

        if err != nil {
          log.Fatal(err)
          break
        }
    }
    if err != nil {
        log.Fatal(err)
        break
    }
  }
}