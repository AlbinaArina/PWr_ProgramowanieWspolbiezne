package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var numRows, numCols int

var matrix [][]int
var mutex = &sync.Mutex{}

func initializeMatrix() {
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			matrix[i][j] = -1
		}
	}
}

func printMatrix() {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("")
	fmt.Println("Rozmieszcenie podróżników:")

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			if matrix[i][j] == -1 {
				fmt.Print("-- ")
			} else if matrix[i][j] == -11 {
			    matrix[i][j] = -1
				fmt.Print("<= ")
			} else if matrix[i][j] == -22 {
			    matrix[i][j] = -1
				fmt.Print("=> ")
			} else if matrix[i][j] == -33 {
			    matrix[i][j] = -1
				fmt.Print("|| ")
			} else {
				fmt.Printf("%02d ", matrix[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
    
	fmt.Print("Podaj liczbę wierszy: ")
	fmt.Scan(&numRows)
	fmt.Print("Podaj liczbę kolumn: ")
	fmt.Scan(&numCols)

	
	matrix = make([][]int, numRows)
    for i := range matrix {
        matrix[i] = make([]int, numCols)
        for j := range matrix[i] {
            matrix[i][j] = -1
        }
    }
    
	rand.Seed(time.Now().UnixNano())
	initializeMatrix()

	go func() {
		currentNumber := 0
	    for {
			// Find a random empty spot
			mutex.Lock()
			emptySpots := make([][2]int, 0)
			for i := 0; i < numRows; i++ {
				for j := 0; j < numCols; j++ {
					if matrix[i][j] == -1 {
						emptySpots = append(emptySpots, [2]int{i, j})
					}
				}
			}

			if len(emptySpots) > 0 {
				// Place the next number in a random empty spot
				index := rand.Intn(len(emptySpots))
				matrix[emptySpots[index][0]][emptySpots[index][1]] = currentNumber
				currentNumber++
			}

			mutex.Unlock()
			time.Sleep(time.Second)
			time.Sleep(time.Second)
		}
		
		
	}()

	go func() {
		for {
			mutex.Lock()
			for i := 0; i < numRows; i++ {
				for j := 0; j < numCols; j++ {
					if matrix[i][j] >= 0 {
						// Choose a random action
						action := rand.Intn(5)
						switch action {
						case 0: // Stay in place
						case 1: // Move left if possible
							if j > 0 && matrix[i][j-1] == -1 {
								matrix[i][j-1] = matrix[i][j]
								matrix[i][j] = -11
							}
						case 2: // Move right if possible
							if j < numCols-1 && matrix[i][j+1] == -1 {
								matrix[i][j+1] = matrix[i][j]
								matrix[i][j] = -22
							}
						case 3: // Move up if possible
							if i > 0 && matrix[i-1][j] == -1 {
								matrix[i-1][j] = matrix[i][j]
								matrix[i][j] = -33
							}
						case 4: // Move down if possible
							if i < numRows-1 && matrix[i+1][j] == -1 {
								matrix[i+1][j] = matrix[i][j]
								matrix[i][j] = -33
							}
						}
					}
				}
			}
			mutex.Unlock()
			printMatrix()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Periodically print the matrix
	for {
		printMatrix()
		time.Sleep(2*time.Second)
	}
}

