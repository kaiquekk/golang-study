package main

import (
	"fmt"
	"sync"
	"time"
)

type Fork struct {
	sync.Mutex
}

type Philosopher struct {
	name                string
	leftFork, rightFork *Fork
}

func (p Philosopher) eat() {
	p.leftFork.Lock()
	p.rightFork.Lock()

	fmt.Printf("%v%v is eating.\n", string("\033[31m"), p.name)
	time.Sleep(5 * time.Second)

	p.leftFork.Unlock()
	p.rightFork.Unlock()

	fmt.Printf("%v%v has finished eating.\n", string("\033[32m"), p.name)
	time.Sleep(1 * time.Second)
}

func main() {
	var wg sync.WaitGroup

	nPhilosophers := 5
	nForks := 5
	philosophers := make([]*Philosopher, nPhilosophers)
	forks := make([]*Fork, nForks)

	for i := 0; i < nPhilosophers; i++ {
		forks[i] = new(Fork)
	}

	for i := 0; i < nPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			name:      fmt.Sprintf("Philosopher %v", i),
			rightFork: forks[i],
			leftFork:  forks[(i+1)%nForks],
		}
	}

	for i := 0; i < nPhilosophers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			philosophers[i].eat()
		}(i)
	}

	wg.Wait()

	fmt.Print(string("\033[0m"))
}
