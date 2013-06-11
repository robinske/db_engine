import socket
import sys
import string
import unicodedata

HOST, PORT = "localhost", 4127
LOGFILE = "../outputs/log.txt"

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect((HOST, PORT))

def connect(DATABASE):
		sock.sendall("DATABASE:> "+DATABASE)
		
		r = sock.recv(4096)
		print "{}".format(r)
		
		apply_log()

		return

def save():
		sock.sendall("SAVE")

def apply_log():
		f = open(LOGFILE)
		
		filetext = f.read()
		
		filelist = filetext.split("\n")

		if len(filelist) > 0:
				for line in filelist:
							if line != "":
									sock.send(line)
									r = sock.recv(1024)
									print "{}".format(r)
		f.close()
		sock.send("SAVE")
		return

def query(instruction):
	
	while True:

		data = str.upper(instruction)
			    
		if data is not None:	    
		    sock.sendall(data)
		    received = sock.recv(4096)
		    return_data = "{}".format(received)

		    return return_data

		if data == "QUIT":
			print "Goodbye!"
			sock.close()
			sys.exit(0)