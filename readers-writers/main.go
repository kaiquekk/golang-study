package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var mutex = new(sync.RWMutex)

var nToBeStored int

func writer(wg *sync.WaitGroup) {
	defer wg.Done()

	var wg2 sync.WaitGroup
	nWriters := 10

	for i := 0; i < nWriters; i++ {
		fmt.Printf("%vStarted Writer num: %v\n", string("\033[36m"), i)
		wg2.Add(1)

		time.Sleep(time.Second)

		go func(i int) {
			defer wg2.Done()

			mutex.Lock()
			defer mutex.Unlock()
			nToBeStored = rand.Intn(100)

			fmt.Printf("%vLocked Writer num: %v, Stored Value: %v\n", string("\033[31m"), i, nToBeStored)
		}(i)
	}

	wg2.Wait()
}

func reader(wg *sync.WaitGroup) {
	defer wg.Done()

	var wg2 sync.WaitGroup
	numReaders := 10

	for i := 0; i < numReaders; i++ {
		fmt.Printf("%vStarted Reader num: %v\n", string("\033[33m"), i)
		wg2.Add(1)

		time.Sleep(time.Second)

		go func(i int) {
			defer wg2.Done()

			mutex.RLock()
			defer mutex.RUnlock()

			fmt.Printf("%vLocked Reader num: %v, Stored Value: %v\n", string("\033[32m"), i, nToBeStored)
		}(i)
	}

	wg2.Wait()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go writer(&wg)
	go reader(&wg)
	wg.Wait()

	fmt.Println(string("\033[0m"))
}
