package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	tobacco = 0
	paper   = 1
	match   = 2
)

var ingredientMap = []string{"tobacco", "paper", "match"}

type Table struct {
	ingredients [2]int
	mutex       sync.Mutex
	condition   *sync.Cond
}

func (t *Table) placeIngredients(ingredientA, ingredientB int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.ingredients[0] = ingredientA
	t.ingredients[1] = ingredientB
	t.condition.Broadcast()
}

func (t *Table) takeIngredients(ingredient int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for t.ingredients[0] == -1 || t.ingredients[1] == -1 {
		t.condition.Wait()
	}
	if t.ingredients[0] != ingredient && t.ingredients[1] != ingredient {
		t.ingredients[0], t.ingredients[1] = -1, -1
		return true
	}
	return false
}

func agent(table *Table) {
	for {
		time.Sleep(3 * time.Second)
		ingredientA := rand.Intn(3)
		ingredientB := rand.Intn(3)
		for ingredientA == ingredientB {
			ingredientB = rand.Intn(3)
		}
		fmt.Printf("%vAgent places %v%s%v and %v%s%v on the table.%v\n",
			string("\033[32m"),
			string("\033[36m"),
			ingredientMap[ingredientA],
			string("\033[32m"),
			string("\033[36m"),
			ingredientMap[ingredientB],
			string("\033[32m"),
			string("\033[0m"),
		)
		table.placeIngredients(ingredientA, ingredientB)
	}
}

func smoker(table *Table, ingredientInPoss int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if table.takeIngredients(ingredientInPoss) {
			fmt.Printf("%vSmoker with %v%s%v makes a cigarette.%v\n", string("\033[31m"), string("\033[36m"), ingredientMap[ingredientInPoss], string("\033[31m"), string("\033[0m"))
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	table := new(Table)
	table.ingredients[0] = -1
	table.ingredients[1] = -1
	table.condition = sync.NewCond(&table.mutex)
	var wg sync.WaitGroup

	wg.Add(3)
	go agent(table)
	go smoker(table, tobacco, &wg)
	go smoker(table, paper, &wg)
	go smoker(table, match, &wg)

	wg.Wait()
}
