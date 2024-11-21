package as2

import (
	"fmt"
	"time"
)

func livelock() {
	var x int
	y := make(chan int, 1)

	// T2
	go func() {
		y <- 1
		x++
		<-y

	}()

	x++
	y <- 1
	<-y

	time.Sleep(1 * 1e9)
	fmt.Printf("done \n")

}

func dl_snd(ch chan int) {
	var x int = 0
	x++
	ch <- x
}

func dl_rcv(ch chan int) {
	var x int
	x = <-ch
	fmt.Printf("received %d \n", x)

}

func deadlock() {
	var ch chan int = make(chan int)
	go dl_rcv(ch) // R1
	go dl_snd(ch) // S1
	dl_rcv(ch)    // R2

}

func datarace(delay int, y chan int) {
	var x int

	// T2
	go func() {
		y <- 1
		x++
		<-y

	}()

	x++
	y <- 1
	<-y

	dl := delay * int(time.Second)
	time.Sleep(time.Duration(dl))
	fmt.Printf("done \n")

}
