// Authors: Evan Menendez, Sarah Haddix, Ben Chesser
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numChairs    = 2
	numCustomers = 100
)

func barber(waitingRoom chan int, barberWg *sync.WaitGroup) {
	for {
		fmt.Println("Barber is sleeping")
		customer, ok := <-waitingRoom
		if !ok {
			fmt.Println("Barber is closing up becuase no customers are waiting")
			return
		}
		fmt.Printf("Barber is cutting hair of customer %d\n", customer)
		time.Sleep(time.Second)
		fmt.Printf("Barber finished cutting hair of customer %d\n", customer)
		barberWg.Done()
	}
}

func customer(id int, waitingRoom chan int, customerWg, barberWg *sync.WaitGroup) {
	defer customerWg.Done()
	fmt.Printf("Customer %d pulls up to barber shop\n", id)
	select {
	case waitingRoom <- id:
		fmt.Printf("Customer %d sits in the waiting room\n", id)
		barberWg.Add(1)
	default:
		fmt.Printf("Customer %d leaves beacuse its busy\n", id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	waitingRoom := make(chan int, numChairs)
	var customerWg sync.WaitGroup
	var barberWg sync.WaitGroup

	go barber(waitingRoom, &barberWg)

	for i := 1; i <= numCustomers; i++ {
		customerWg.Add(1)
		go customer(i, waitingRoom, &customerWg, &barberWg)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	customerWg.Wait()
	close(waitingRoom)
	barberWg.Wait()
}
