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

func talkToDB() (cacheData) {
  data := make(cacheData)
  fmt.Println("YOU'VE TALKED TO THE DB")
  return data
}

func get(key, instruct string) {
  fmt.Printf("You called the %s function\n", instruct)
  fmt.Printf("We'll get the value of key %s\n", key)
  
  //fmt.Println(data[key])
}

func put(key, value string) {
  // fmt.Printf("You called the %s function\n", instruct)
  fmt.Printf("We'll put %s:%s in the dictionary\n", key, value)
  //data := cacheData{key:value}
  //fmt.Println(data[key]) // should print value
}

func save(key, value, instruct string) {
  fmt.Printf("You called the %s function\n", instruct)
  fmt.Printf("We'll save %s:%s from dictionary to disk\n", key, value)
}

func parseRequest(message string) (string) {
  msgSplit := strings.Split(message, " ")

  //fmt.Printf("%s", message)
  //fmt.Printf("%s", msgSplit)

  instruct := msgSplit[0]
  key := msgSplit[1]
  value := msgSplit[2]

  // if len(msgSplit) < 2 {
    // handle error
  // }
  // database := "database" // could be a "FROM" statement later

  

  switch instruct {
    case "GET": get(key, instruct)
    case "PUT": put(key, value)
    case "SAVE": save(key, value, instruct)
    default: fmt.Println("try again idiot")
  }

  return value
}

func addRequestToDict(message string) {


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