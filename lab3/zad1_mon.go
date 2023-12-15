package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhilosophers = 5

type PhilosopherMonitor struct {
	forks     [numPhilosophers]int
	okToEat   [numPhilosophers]*sync.Cond
	philMutex sync.Mutex
}

func NewPhilosopherMonitor() *PhilosopherMonitor {
	pm := &PhilosopherMonitor{}
	for i := 0; i < numPhilosophers; i++ {
		pm.forks[i] = 2
		pm.okToEat[i] = sync.NewCond(&pm.philMutex)
	}
	return pm
}

func (pm *PhilosopherMonitor) takeFork(i int) {
	pm.philMutex.Lock()
	for pm.forks[i] != 2 {
		pm.okToEat[i].Wait()
	}
	pm.forks[(i+1)%numPhilosophers]--
	pm.forks[(i-1+numPhilosophers)%numPhilosophers]--
	pm.philMutex.Unlock()
}

func (pm *PhilosopherMonitor) releaseFork(i int) {
	pm.philMutex.Lock()
	pm.forks[(i+1)%numPhilosophers]++
	pm.forks[(i-1+numPhilosophers)%numPhilosophers]++
	if pm.forks[(i+1)%numPhilosophers] == 2 {
		pm.okToEat[(i+1)%numPhilosophers].Signal()
	}
	if pm.forks[(i-1+numPhilosophers)%numPhilosophers] == 2 {
		pm.okToEat[(i-1+numPhilosophers)%numPhilosophers].Signal()
	}
	pm.philMutex.Unlock()
}

func philosopher(id int, pm *PhilosopherMonitor, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		think(id)
		pm.takeFork(id)
		eat(id)
		pm.releaseFork(id)
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

	// Create philosopher monitor
	pm := NewPhilosopherMonitor()

	// Create philosophers
	for i := 0; i < numPhilosophers; i++ {
		wg.Add(1)
		go philosopher(i, pm, &wg)
	}

	wg.Wait()
}