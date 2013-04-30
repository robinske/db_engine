func update(connection net.Conn, key, value string) {

	keyList := strings.Split(key, "/")
	collection := keyList[0]
	var index int

	for k, v := range lock.nested {
		if k == collection {
			// if v is an array:
			for i, element := range v {
				fmtElem := fmt.Sprintf("%v", element)
				if strings.Contains(fmtElem, keyList[2]) {
					index = i
				}
			}
			v[index][keyList[3]] = keyList[4] // this will give you the list element you need to append
			connection.Write([]byte("We updated something"))
		} else {
			connection.Write([]byte("nothing found to update//error"))
		}

			// value := ksearch(connection, keyList[1], keyList[2])
			// find the index of v where that value lives
			// v[keyList[3]]
			// use the value to search for that unique key

	}

}


given an input: update users/id/941/zip 49684
			key = users/id/941/zip
			new value = 49684
			

