package race

// Code from https://github.com/aditya43/golang_concurrency/blob/main/05-race-detector/main.go

import (
	"fmt"
	"math/rand"
	"time"
)

// identify the data race
// fix the issue.

func race(maxDuration int) {
	start := time.Now()
	reset := make(chan bool)
	// var t *time.Timer
	t := time.AfterFunc(randomDuration(maxDuration), func() {
		fmt.Println(time.Since(start))
		reset <- true
	})
	for time.Since(start) < 5*time.Second {
		<-reset
		t.Reset(randomDuration(maxDuration))
	}
}

func randomDuration(maxDuration int) time.Duration {
	maxTime := maxDuration * int(time.Second)
	return time.Duration(rand.Int63n(int64(maxTime)))
}

//----------------------------------------------------
// (main goroutine) -> t <- (time.AfterFunc goroutine)
//----------------------------------------------------
// (working condition)
// main goroutine..
// t = time.AfterFunc()  // returns a timer..

// AfterFunc goroutine
// t.Reset()        // timer reset
//----------------------------------------------------
// (race condition- random duration is very small)
// AfterFunc goroutine
// t.Reset() // t = nil

// main goroutine..
// t = time.AfterFunc()
//----------------------------------------------------
