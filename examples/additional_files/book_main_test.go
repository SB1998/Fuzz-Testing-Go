package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	// Create a Wait Group, and save the pointer
	wg := &sync.WaitGroup{}

	// Create a Mutex, once against using a pointer
	m := &sync.RWMutex{}

	cacheCh := make(chan Book)
	dbCh := make(chan Book)

	for i := 0; i < 10; i++ {
		// Fetch a random Book ID
		id := rnd.Intn(10) + 1

		// We've got 2 GoRoutines to wait for
		wg.Add(2)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			// Query the cache, if found then print it out
			if b, ok := queryCache(id, m); ok {
				ch <- b // Pass the found Book to the cache channel
			}

			wg.Done()
		}(id, wg, m, cacheCh)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			// Query the DB, if found then print it out
			if b, ok := queryDatabase(id, m); ok {
				m.Lock()
				cache[id] = b
				m.Unlock()
				ch <- b
			}

			wg.Done()
		}(id, wg, m, dbCh)

		// Create on GoRoutine per query to handle response
		go func(cacheCh, dbCh <-chan Book) {
			select {
			case b := <-cacheCh:
				fmt.Println("Source: Cache")
				fmt.Println(b)
				<-dbCh // Wait to get the message from the DB Channel so we don't block
			case b := <-dbCh:
				fmt.Println("Source: Database")
				fmt.Println(b)
			}
		}(cacheCh, dbCh)

		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait()
}
