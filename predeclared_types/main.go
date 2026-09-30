package main

import "fmt"

func main() {
	var myFirstInitial rune = 'J'
	var myLastInitial int32 = 'J'

	fmt.Println("first: ", myFirstInitial)
	fmt.Println("last: ", myLastInitial)

	// need to use type conversion when variable types do not match
	var x int = 10
	var y float64 = 30.2
	var sum1 float64 = float64(x) + y
	var sum2 int = x + int(y)
	fmt.Println(sum1, sum2)
}
