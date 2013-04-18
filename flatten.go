package main

type Dict map[string]string
type Dict2 map[string]interface{}

type Tuple struct {
	nk string
	val string
}

var parentKey = ""

func flatten(dictionary Dict, parentKey string) Dict {
	var dictItems []string
	for k, v := range dictionary {
		if parentKey != "" {
			newKey := parentKey + "_" + k	
		} else {
			newKey := k
		}
		// if isinstance(v, collections.MutableMapping) NEED THIS EQUIVALENT IN GOLANG (tried below)
		var ok bool
		if _, ok := v.(Dict); ok { // this should check that value is of type Dict 
			dictItems = append(dictItems, flatten(v, newKey)) // this probably doesn't work either 
		} else {
			var t = Tuple {}
			t.nk = newKey
			t.val = v
			dictItems = append(dictItems, t) // so apparently Go doesn't support tuples
		}
	}
	return dict(items) // need golang equivalent
}


func main() {

	var d = Dict {"key": "value"}

	flatten(d)

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