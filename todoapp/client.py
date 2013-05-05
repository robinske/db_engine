import socket
import sys
import string
import unicodedata

HOST, PORT = "localhost", 4127

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect((HOST, PORT))

def connect(DATABASE):
	sock.sendall("DATABASE:> "+DATABASE)
	r = sock.recv(1024)
	print "{}".format(r)

def save():
	sock.sendall("SAVE")

def apply_log():

#     fileContents, err := ioutil.ReadFile(LOGFILE)
#     fileString := string(fileContents)

#     fileArray := strings.Split(fileString, "\n")
#     fileArray = append(fileArray, "CLEARLOG")

#     for _, line := range fileArray {
#         if line != "" {
#             connection.Write([]byte(line))
#             inputEnd, err := connection.Read(buf[:])
#             if err != nil {
#                 return
#             }
#             fmt.Printf("%s\n", string(buf[0:inputEnd]))
#         }
#     }
# }

def is_log_empty():

#     fileContents, err := ioutil.ReadFile(LOGFILE)
#     fileString := string(fileContents)

#     if err != nil {
#         log.Fatal(err)
#     }

#     if fileString == "" {
#         return true
#     }
#     return false
# }

def query(instruction):
	
	while True:

		data = str.upper(instruction)
			    
		if data is not None:	    
		    sock.sendall(data)
		    received = sock.recv(4096)
		    return_data = "{}".format(received)

		    print return_data
		    return return_data

		if data == "QUIT":
			print "Goodbye!"
			sock.close()
			sys.exit(0)