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