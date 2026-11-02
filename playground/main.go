package main

import (
	"fmt"
	"os"
)

func run() int {
	defer fmt.Println("cleanup ran")
	return 1
}

func main() {
	// convert slice of ints into comma separated string
	// nums := []int64{1,2,3,4,5,6}
	// result := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(nums)), ","), "[]")
	// fmt.Println(result)

	// try each line at a time and see what happens
	os.Exit(run())
	defer fmt.Println("cleanup ran")
	panic("boom")
	os.Exit(256)
}
