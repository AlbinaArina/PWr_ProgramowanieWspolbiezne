package main

import (
	"fmt"
	"sync"
	"time"
)

type RWMonitor struct {
	mu          sync.Mutex
	okToRead    *sync.Cond
	okToWrite   *sync.Cond
	readers     int
	writing     bool
}

func NewRWMonitor() *RWMonitor {
	rwm := &RWMonitor{}
	rwm.okToRead = sync.NewCond(&rwm.mu)
	rwm.okToWrite = sync.NewCond(&rwm.mu)
	return rwm
}

func (rwm *RWMonitor) startRead() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	for rwm.writing {
		rwm.okToRead.Wait()
	}

	rwm.readers++
	fmt.Println("Reader started reading. Current amount of readers:", rwm.readers)
	rwm.okToRead.Signal()
}

func (rwm *RWMonitor) stopRead() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	rwm.readers--
	if rwm.readers == 0 {
		rwm.okToWrite.Signal()
	} else {
		fmt.Println("Reader finished reading. Current amount of readers:", rwm.readers)
	}
}

func (rwm *RWMonitor) startWrite() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	for rwm.readers > 0 || rwm.writing {
		rwm.okToWrite.Wait()
	}

	rwm.writing = true
	fmt.Println("Writer started writing.")
}

func (rwm *RWMonitor) stopWrite() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	rwm.writing = false
	if rwm.readers > 0 {
		rwm.okToRead.Signal()
	} else {
		rwm.okToWrite.Signal()
	}
	fmt.Println("Writer finished writing.")
}

func reader(rwm *RWMonitor, id int) {
	for {
		time.Sleep(time.Second)
		rwm.startRead()
		// Czytaj dane
		fmt.Printf("Reader %d is reading.\n", id)
		rwm.stopRead()
	}
}

func writer(rwm *RWMonitor, id int) {
	for {
		time.Sleep(2 * time.Second)
		rwm.startWrite()
		// Zapisz dane
		fmt.Printf("Writer %d is writing.\n", id)
		rwm.stopWrite()
	}
}

func main() {
	rwm := NewRWMonitor()

	for i := 1; i <= 6; i++ {
		go reader(rwm, i)
	}

	for i := 1; i <= 3; i++ {
		go writer(rwm, i)
	}

	// Program działa przez pewien czas
	time.Sleep(10 * time.Second)
}

