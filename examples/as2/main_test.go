package as2

import (
	"fmt"
	"testing"
	"time"
)

func TestLiveLock(t *testing.T) {
	livelock()

}

func TestDeadLock(t *testing.T) {
	deadlock()

}

func TestDataRace(t *testing.T) {
	y := make(chan int, 1)
	datarace(1, y)

	select {
	case <-time.After(5 * time.Second):
		fmt.Printf("Timeout!")
	case e := <-y:
		fmt.Printf("Received %+v", e)
	}

}
