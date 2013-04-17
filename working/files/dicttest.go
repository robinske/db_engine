package main

import "fmt"

type friends map[string]int

func main() {
	gang := friends {}
	
	fmt.Println(gang)

	gang["divya"] = 1
	gang["kelley"] = 2

	fmt.Println(gang)
}