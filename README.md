Database Engine
=====================

Overview
-------------
GoSQL is an in-memory data store in the style of Redis or MongoDB. Data remains in memory for the duration of execution while simultaneously being persisted to disk.
 
Despite the name, there is no SQL in GoSQL.

Languages:

    * Go (server and client)
    * Python (client for web application)

Network Protocol
------------------------
Data is sent over a TCP/IP socket as a raw byte stream. The client sends a single command with parameters separated by spaces followed by a newline.

The server responds with a chunk of JSON data. The length of the message is not transmitted prior to this, which means the client simply reads until there is no more message data left. In theory, there are cases where the client will time out waiting for data. In practice, the end of a message is determined when there is less data than expected, so the chance of a timeout is very low: it only happens when there is exactly as much data as expected.

The server and client communicate by default on port 4127. Implemented using the Go Net package, the server waits for a client connection to communicate using a Go byte array.

Storage Schema
------------------------

Data is stored to and loaded from disk as JSON. On load, data is unmarshaled from any valid JSON string using the Go JSON package. Keys must be strings but values can be any type, including other JSON objects or arrays. Data is repackaged as JSON (line 418, line 573) when saved to disk.

In-Memory Data Storage
-----------------------------------

The in-memory data storage is a nested Go hash map, where the keys are strings and values are an empty interface type. Because of this, when data is loaded, it is inserted into cache storage by recursively type checking values in order to account for different nested data types (lines 327 - 385).

Resilience Mechanism
-------------------------------

The database is written in entirety to disk as JSON on clean shutdown (given the command "quit"). Throughout the rest of the operations, any changes to the data structure are logged to disk. On startup, the most recent full save is loaded, then the log is applied to recreate the last known state of the data store. 

Concurrency
------------------

Considerations were implemented by using's Go's built in concurrency model ("goroutine") as seen on line 603. This allows any number of concurrent clients to interact with the server.

In order to preserve the in memory data store for concurrent use, a semaphore/mutex is applied for both read and write access during any database read or update. Cache object is initiated on lines 17 - 20.

The mutex ensures consistency of data across multiple access points, but does create a queue as only one client can read or edit in memory data at any given time. Removing the read locks may decrease consistency but would increase throughput if speed was prioritized.

Server
---------

The server reads and processes information sent to it from any number of concurrent clients. The main loop of the server is seen on line 596. The server infinitely loops to accept incoming messages from the client.

Client
--------

Implemented two clients, one in Python and one in Go. Both have command line interfaces that accept the database query language. Also built a version of the Python client to interact with a Flask web application as the data store.

The interface uses Standard I/O to read incoming messages from the command line (line 70) and implements a buffer to transmit information.

Query Language // Available Operations
----------------------------------------------------------

Because data is backed by a hashmap, most operations mimic the runtime complexity of hashmaps.

  {
      "TASKID": 3,
      "TASKS" : 
                [
                    {
                        "ID": 1,
                        "TITLE": "CHECK EMAIL",
                        "COMPLETED": "NO"
                    },
                    {
                        "ID": 2,
                        "TITLE": "DO DISHES",
                        "COMPLETED": "NO"
                    }
                ]
  }

Valid keys
    * Any string value that does not contain spaces or colons
Valid partial keys
    * Any string value that does not contain spaces
    * Will match any key that contains the partial key as a substring
Valid values
    * Defaults to string, may contain spaces. 
    * May also be any valid JSON value (integers, Booleans, null, arrays, nested objects)
    * Updates to nested values assume it is an array of objects
Note: string data is normalized to ALL CAPS for processing/comparisons.

	  SHOW [COLLECTIONS|DATA|DATABASE]
		    Displays either:
		    - top level keys/collections 
		    - snapshot of current data
		    - database file
		    - ex: ("SHOW COLLECTIONS" displays current schema)

	  GET key
        Returns the value of top level key
        - ex: "GET TASKS" returns the JSON array value of tasks

	  SEARCH partial key
        Displays value of nested keys
        - flattens keys (line 387) into their respective paths, searches resulting keys for nested keys
        - ex: "SEARCH ID" returns a list of IDs:
							"TASKS:1:ID: 2"
							"TASKS:0:ID: 1"

	  SEARCHBYKEY key value
				Displays all keys where, given a key and value, that value is true
				- ex: "SEARCHBYKEY ID 1" returns:
				      "TASKS:0:TASK_ID: 1"

	  SET key value
	      Adds new top level key/collection
	      - ex: "SET NEWKEY NEWVALUE"

	  UPDATE key value
	      Updates top level key/collection (default string values)
	      - ex: "UPDATE TASKID 4" changes TASKID from integer 3 to string 4

	  UPDATEINT key value
	      Updates top level key/collection where value will become an integer
	      - same as update, but "UPDATE TASKID 4" changes TASKID from integer 3 to integer 4

	  NESTEDUPDATE top_level_key/unique_id/key_to_update value
	      Updates a value nested within a JSON object inside an array. It uses the unique id to find the array index and constructs the full key path for the update. Update is then performed like standard UPDATE.
	      - ex: "NESTEDUPDATE TASKS/2/COMPLETED YES" will update that we have completed doing the dishes.
	      - Must provide a unique id for the path, in this case "2"

	  ADDTO key value
	      Set for nested values. Appends new values to the end of a JSON array and reinserts into the cache
	      - ex: "ADDTO TASKS {"TASK_ID":3, "TITLE":"NEW TASK", "COMPLETED":"NO"}" adds a third task to our task list

	  REMOVE key
	      Deletes top level key/collection
	      - ex: "REMOVE NEWKEY" deletes NEWKEY and its data
