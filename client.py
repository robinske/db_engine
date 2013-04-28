import socket
import sys
import string

HOST, PORT = "localhost", 4127

# Create a socket (SOCK_STREAM means a TCP socket)
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect((HOST, PORT))

if len(sys.argv) < 2:
	print "Please load a database"
else:
	DATABASE = sys.argv[1]
	sock.sendall("DATABASE:> "+DATABASE)
	r = sock.recv(1024)
	print "{}".format(r)

while True:
	data = str.upper(raw_input(">>> "))

	if data is not None:	    
	    sock.sendall(data)
	    received = sock.recv(1024)

	print "{}".format(received)

	if data == "QUIT":
		print "Goodbye!"
		sock.close()
		sys.exit(0)