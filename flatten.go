	for rkey, value := range inputJSON {
package main

import (
	"fmt"
	"flatten"
)

type Dict2 map[string]string
type Dict map[string]interface{}

type Tuple struct {
	nk string
	val string
}

// var parentKey = ""

// func flatten(dictionary Dict, parentKey string) Dict {
// 	var dictItems []string
// 	for k, v := range dictionary {
// 		if parentKey != "" {
// 			newKey := parentKey + "_" + k	
// 		} else {
// 			newKey := k
// 		}
// 		// if isinstance(v, collections.MutableMapping) NEED THIS EQUIVALENT IN GOLANG (tried below)
// 		var ok bool
// 		if _, ok := v.(Dict); ok { // this should check that value is of type Dict 
// 			dictItems = append(dictItems, flatten(v, newKey)) // this probably doesn't work either 
// 		} else {
// 			var t = Tuple {}
// 			t.nk = newKey
// 			t.val = v
// 			dictItems = append(dictItems, t) // so apparently Go doesn't support tuples
// 		}
// 	}
// 	return dict(items) // need golang equivalent
// }

// Template data
// data := struct {
// 	Title string
// 	Users []*User
// }{
// 	title,
// 	users,
// }
// err := tmpl.Execute(w, data)
// (Cheaper and safer than using map[string]interface{}.)


var lkey = ""

func flatten(inputJSON Dict, lkey string) map[string]string {
	flattenedJSON := map[string]string {}
	for rkey, value := range inputJSON {
		key := lkey+rkey
		if _, ok := value.(Dict); ok { // type assertion

			// if isinstance(val, dict):
   				// ret.update(parse_dict(val, key+'_'))

			for key, value := range flatten(inputJSON, lkey+"_") {
				flattenedJSON[key] = value.(string)
			}
		} else {
			flattenedJSON[key] = value.(string)
		}
	}
	return flattenedJSON
}

func main() {

	var d = Dict{"key": "value", "key2": {"key3": "value2"}}
	//var d = Dict {"ay": "1", "cey": {"ay": "2", "b": {"x": "5", "y": : "10"}}, "d": "[1, 2, 3]"}

	fmt.Printf("%v\n", flatten(d, lkey))

}

// WORKING PYTHON CODE

// import collections

// def flatten(d, parent_key=''):
//     items = []
//     for k, v in d.items():
//         new_key = parent_key + '_' + k if parent_key else k
//         if isinstance(v, collections.MutableMapping):
//             items.extend(flatten(v, new_key).items())
//         else:
//             items.append((new_key, v))
//     # print items
//     return dict(items)

// print flatten({'a': 1, 'c': {'a': 2, 'b': {'x': 5, 'y' : 10}}, 'd': [1, 2, 3]})