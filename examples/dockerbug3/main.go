package docker

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// this if a coppy of https://github.com/MrDKOz/golang-concurrency
type Entry struct {
	ID            int
	Title         string
	Author        string
	YearPublished int
}

// emulate an object
type Daemon struct{}

var entries = []Entry{
	{
		ID:            1,
		Title:         "The Hitchhiker's Guide to the Galaxy",
		Author:        "Douglas Adams",
		YearPublished: 1979,
	},
	{
		ID:            2,
		Title:         "The Hobbit",
		Author:        "J.R.R Tolkein",
		YearPublished: 1937,
	},
	{
		ID:            3,
		Title:         "A Tale of Two Cities",
		Author:        "Charles Dickins",
		YearPublished: 1859,
	},
	{
		ID:            4,
		Title:         "Harry Potter and the Philosophers Stone",
		Author:        "J.K. Rowling",
		YearPublished: 1997,
	},
	{
		ID:            5,
		Title:         "Les Miserables",
		Author:        "Victor Hugo",
		YearPublished: 1862,
	},
	{
		ID:            6,
		Title:         "I, Robot",
		Author:        "Isaac Asamov",
		YearPublished: 1950,
	},
	{
		ID:            7,
		Title:         "The Gods Themselves",
		Author:        "Isaac Asamov",
		YearPublished: 1973,
	},
	{
		ID:            8,
		Title:         "The Moond is a Hash Mistress",
		Author:        "Robert A. Heinlein",
		YearPublished: 1966,
	},
	{
		ID:            9,
		Title:         "On Basilisk Station",
		Author:        "David Weber",
		YearPublished: 1993,
	},
	{
		ID:            10,
		Title:         "The Android's Dream",
		Author:        "John Scalzi",
		YearPublished: 2006,
	},
}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
var daemon = Daemon{}

// this is the changed implementation of the docker example by gFuzz

func parent() { // parent goroutine
	ch, errCh := daemon.Watch()
	select {
	case <-time.After(1 * time.Nanosecond):
		fmt.Printf("Timeout!")
	case e := <-ch:
		fmt.Printf("Received %+v", e)
	case e := <-errCh:
		fmt.Printf("Error %s", e)
	}
	return
}

func parentFixed() { // parent goroutine
	ch, errCh := daemon.WatchFixed()
	select {
	case <-time.After(1 * time.Nanosecond):
		fmt.Printf("Timeout!")
	case e := <-ch:
		fmt.Printf("Received %+v", e)
	case e := <-errCh:
		fmt.Printf("Error %s", e)
	}
	return
}

func (d *Daemon) Watch() (chan Entry, chan error) {
	ch := make(chan Entry)
	errCh := make(chan error)
	//+ ch := make(chan discovery.Entries, 1)
	//+ errCh := make(chan error, 1)
	go func() { // child goroutine
		id := rnd.Intn(10) + 1
		entries, err := fetch(id)
		if err != nil {
			errCh <- err
		} else {
			ch <- entries
		}
	}()
	return ch, errCh
}

func (d *Daemon) WatchFixed() (chan Entry, chan error) {
	//- ch := make(chan discovery.Entries)
	//- errCh := make(chan error)
	ch := make(chan Entry, 1)
	errCh := make(chan error, 1)
	go func() { // child goroutine
		id := rnd.Intn(10) + 1
		entries, err := fetch(id)
		if err != nil {
			errCh <- err
		} else {
			ch <- entries
		}
	}()
	return ch, errCh
}

func fetch(id int) (Entry, error) {
	for _, b := range entries {
		if b.ID == id {
			return b, nil
		}
	}

	return Entry{}, errors.New("NO SUCH BOOK FOUND")

}
