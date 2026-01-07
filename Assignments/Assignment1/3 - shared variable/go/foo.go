// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	// "time"
)

var i = 0

func incrementing(increment chan int, finished chan int) {
	for j := 0; j < 5; j++ {
		increment <- 1
	}
	finished <- 1
}

func decrementing(decrement chan int, finished chan int) {
	for j := 0; j < 3; j++ {
		decrement <- 1
	}
	finished <- 1
}

func server(increment chan int, decrement chan int, get chan int, finished chan int) {
	for {
		select {
		case <-increment:
			Println("++")
			i++
		case <-decrement:
			Println("--")
			i--
		case <-get:
			Println("The magic number is:", i)
			finished <- 1
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)

	increment := make(chan int)
	decrement := make(chan int)
	get := make(chan int)

	finished := make(chan int)

	go server(increment, decrement, get, finished)
	go incrementing(increment, finished)
	go decrementing(decrement, finished)

	<-finished
	<-finished
	get <- 1
	<-finished
}

//
