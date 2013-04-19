package flat

// type JSON map[string]interface{}

// var lkey = ""

func flatten(inputJSON map[string]interface{}, lkey string) map[string]string {
	flattenedJSON := map[string]string {}
	for rkey, value := range inputJSON {
		key := lkey+rkey
		if _, ok := value.(map[string]interface{}); ok { // type assertion

			// if isinstance(val, dict):
   				// ret.update(parse_dict(val, key+'_'))

			for key, value := range flatten(inputJSON, lkey+"_") {
				flattenedJSON[key] = value
			}
		} else {
			flattenedJSON[key] = value.(string)
		}
	}
	return flattenedJSON
}
