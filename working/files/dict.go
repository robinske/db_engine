package main

import (
	"fmt"
	//"sort"
)


func main() {
	x := make(map[string]int)
	x["key"] = 10
	fmt.Println(x["key"])

	y := [6]string{"a","b","c","d","e","f"}
	fmt.Println(y[2:5])

	z := make([]int, 3, 9)
	fmt.Println(len(z))
	//fmt.Println(z)

	a := []int{
	    48,96,86,68,
	    57, 1,63,70,
	    37,34,83,27,
	    19,97, 9,17,
	}
	//fmt.Println(a)

	var lowest = a[i]

	for i := 0; i < len(a); i++ {
		if a[i] < lowest {
			lowest = a[i]
		}
	}

	fmt.Println(lowest)
	// dict := make(map[string]int)

	// dict["user"] = "kelley"
	// dict["pass"] = "password"
	// dict["testint"] = 5

	// fmt.Printf(dict["testint"])

}



