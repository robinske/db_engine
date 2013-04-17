package main

import (
	"fmt"
	//"code.google.com/p/go-tour/wc"
	"math"
	"runtime"
	"time"
	"net/http"
	"os"
	"log"
)

func mult(x, y int) (int) {
	result := x*y
	return result
}

func printSixes() {
	for i := 0; i < 30; i++ {
		if i%6==0 {
			print(i,"\n")
		}
	}
}

func structs() {
	type Hi struct {
		X int
		Y int
	}
	p := Hi{1, 2}
	q := &p
	q.X = 1e9
	fmt.Println(p)
}

func newList() {
	newList := []int{2,3,4,5,6,7,8}
	fmt.Println("newList ==", newList)

	for i := 0; i < len(newList); i++ {
		fmt.Printf("newList[%d] == %d\n", i, newList[i])
	}
}

func WordCount (s string) map[string]int {
	return map[string]int{"x": 1}
}

func hypot () {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

	fmt.Println(hypot(3,4))
}

func adder() func(int) int {
	sum := 0
	return func(x int) int{
		sum += x
		return sum
	}
}

func posNeg() {
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}
}

func fibonacci() func() int {
	x := 0
	y := 1
	return func() int {
		x,y = y,x+y
		return x
	}
}

func runFib() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

func showSys() {
	fmt.Print("Go runs on ")
    switch os := runtime.GOOS; os {
    case "darwin":
        fmt.Println("OS X.")
    case "linux":
        fmt.Println("Linux.")
    default:
        // freebsd, openbsd,
        // plan9, windows...
        fmt.Printf("%s.", os)
    }
}

func Saturday() {
	fmt.Println("When's Friday?")
	today := time.Now().Weekday()
	switch time.Friday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}
}

func printTime() {
	t := time.Now()
	fmt.Println(time.Now())
	switch {
	case t.Hour() < 12:
		fmt.Println("Good Morning!")
	case t.Hour() < 17:
		fmt.Println("Good Afternoon!")
	default:
		fmt.Println("Good Evening!")
	}	
}

type Hello struct{}

func (h Hello) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func runOnWeb() {
	var h Hello
	http.ListenAndServe("localhost:4000", h)
}

func main() {
	file, err := os.Open("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes: %q\n", count, data[:count])
}
























