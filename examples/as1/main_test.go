package as1

import (
	"fmt"
	"testing"
	"time"
)

func TestPhilo1(t *testing.T) {
	var forks = make(chan int, 3)
	forks <- 1
	forks <- 1
	forks <- 1
	go philo1(1, forks) // P1
	go philo1(2, forks) // P2
	philo1(3, forks)    // P3
}

func TestPhilo1_2(t *testing.T) {
	var forks = make(chan int, 3)
	forks <- 1
	forks <- 1
	forks <- 1

	// Run philosophers with a timeout to detect deadlocks
	done := make(chan bool)
	go func() {
		go philo1(1, forks) // P1
		go philo1(2, forks) // P2
		philo1(3, forks)    // P3
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("Test completed without deadlock")
	case <-time.After(5 * time.Second):
		t.Error("Test timed out, possible deadlock detected")
	}
}

func TestPhilo2(t *testing.T) {
	var forks = make(chan int, 3)
	forks <- 1
	forks <- 1
	forks <- 1
	go philo2(1, forks)
	go philo2(2, forks)
	philo2(3, forks)
}

func TestPhilo3(t *testing.T) {
	var forks = make(chan int, 3)
	forks <- 1
	forks <- 1
	forks <- 1
	go philo3(1, forks) // P1
	go philo3(2, forks) // P2
	philo3(3, forks)    // P3
}
