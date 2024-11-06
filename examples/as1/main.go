package as1

import (
	"fmt"
	"time"
)

func philo1(id int, forks chan int) {

	for {
		<-forks
		<-forks
		fmt.Printf("%d eats \n", id)
		time.Sleep(1 * 1e9)
		forks <- 1
		forks <- 1

		time.Sleep(1 * 1e9) // think

	}

}

func philo2(id int, forks chan int) {
	for {
		<-forks
		select {
		case <-forks:
			fmt.Printf("%d eats \n", id)
			time.Sleep(1 * 1e9)
			forks <- 1
			forks <- 1

			time.Sleep(1 * 1e9) // think
		default:
			forks <- 1
		}
	}

}

func philo3(id int, forks chan int) {
	for {
		<-forks
		select {
		case <-forks:
			fmt.Printf("%d eats \n", id)
			time.Sleep(1 * 1e9)
			forks <- 1
			forks <- 1

			time.Sleep(1 * 1e9) // think
		default:
			// forks <- 1  // (LOC)
		}
	}

}
