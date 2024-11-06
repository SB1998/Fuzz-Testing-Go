package main

import (
	"fmt"
	"testing"

	"github.com/dsinecos/go-misc-patterns/dice"
)

func TestDice(t *testing.T) {
	diceValue := make(chan int)

	dice.Throw(diceValue)

	select {
	case a := <-diceValue:
		fmt.Printf("Dice rolled to %v \n", a)
	}
}
