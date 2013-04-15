// SERVER LISTENS
// http://stackoverflow.com/questions/2886719/unix-sockets-in-go

package main
import (
    "net"
    // "net/http"
    "log"
    "io"
    "os"
    // "fmt"
)

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