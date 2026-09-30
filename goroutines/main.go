package main

import (
	"fmt"
	"time"
)

func worker(id int, sec time.Duration, ch chan string) {
	time.Sleep(sec)
	ch <- fmt.Sprintf("worker %d finished", id)
}

func main() {
	ch := make(chan int, 2)

	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
	}()

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
