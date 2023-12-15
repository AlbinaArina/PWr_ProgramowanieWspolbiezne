package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhilosophers = 5

var (
	room      = make(chan struct{}, 4) // Semaphore (Room)
	forks     [numPhilosophers]chan struct{}
	philMutex sync.Mutex
)

func philosopher(id int, leftFork, rightFork chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		think(id)
		room <- struct{}{} // Wait(Room)
		<-leftFork         // Wait(Fork(I))
		<-rightFork        // Wait(Fork((I+1) mod 5))
		eat(id)
		leftFork <- struct{}{}  // Signal(Fork(I))
		rightFork <- struct{}{} // Signal(Fork((I+1) mod 5))
		<-room                   // Signal(Room)
	}
}

func think(id int) {
	fmt.Printf("Philosopher %d is thinking\n", id)
	time.Sleep(time.Millisecond * 500)
}

func eat(id int) {
	fmt.Printf("Philosopher %d is eating\n", id)
	time.Sleep(time.Millisecond * 500)
}

func main() {
	var wg sync.WaitGroup

	// Initialize forks
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = make(chan struct{}, 1) // Binary Semaphore (Fork)
		forks[i] <- struct{}{}              // Initialize forks as available
	}

	// Create philosophers
	for i := 0; i < numPhilosophers; i++ {
		wg.Add(1)
		go philosopher(i, forks[i], forks[(i+1)%numPhilosophers], &wg)
	}

	wg.Wait()
}
