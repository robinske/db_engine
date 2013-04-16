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

func client() {
  c, err := net.Dial("tcp", ":4127") // sets a connection, c, to the port 4127
  if err != nil {
    log.Fatal(err)
  }
  defer c.Close()

  go reader(c) // concurrent process / Goroutine
  
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

func server() {
  l, err := net.Listen("tcp", ":4127")
  if err != nil {
    log.Fatal(err)
    return
  }

  for {
    fd, err := l.Accept()
    if err != nil {
        log.Fatal(err)
        return
    }

    go echoServer(fd)
  }
}

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

  // need to handle message length error - only works for 3+ word inputs right now

  var dictionary = make(cacheData)
  fmt.Println("printing base dictionary", dictionary)

  talkToDictionary(instruct, key, value, dictionary)

  return value, dictionary
}

///

func echoServer(c net.Conn) {
  for {
    // makes a list of 512 byte elements / why 512? arbitrary? cuts off messages after 512 bytes, sends over in a different package size.
    // makes a buffer to keep chunks that are read/written
    buf := make([]byte, 512)
    // sets two variables: nr (stands for??) and err to read the byte elements
    nr, err := c.Read(buf)
    // ignore errors that aren't nil
    if err == io.EOF {
        return
    }
    // sets a variable of "data" to a slice of buf
    data := buf[0:nr]
    
// WRITE TO DATABASE / ONLY CALL THIS IF SAVE?

    fo, err := os.OpenFile("output", os.O_RDWR|os.O_APPEND, 0666) // 0666 is the tag for who can read and write to the file per system reqs
    fo.Seek(0,2) // 2 means go to the end of the file, 0 is the relative position to the end
    if err != nil {
      log.Fatal(err)
    }
    
    defer fo.Close()

    _, err = fo.Write(data) // write to a file!!! / Make this optional file input string

// WRITE TO DATABASE END

    println("Server received:", string(data)) // have it store this to a file

    _, err = c.Write(data)
    if err != nil {
      log.Fatal(err)
    }
  }
}

func main() {
  go server()
  go client()
}