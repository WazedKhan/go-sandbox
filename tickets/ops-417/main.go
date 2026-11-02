package main

// Title: Checkout freezes when cart has more than one item

// Description:
// Customers with a single item check out fine. As soon as a cart has 2+ items, the whole thing freezes and the service dies.
// QA says “it worked in testing” (they only tested with one item 🙃). Sales is yelling at me. Please fix ASAP.

import (
	"fmt"
	"sync"
)

type Inventory struct {
	mu    sync.Mutex
	stock map[string]int
}

func (inv *Inventory) Reserve(items []string) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	for _, item := range items {
		if inv.stock[item] == 0 {
			return fmt.Errorf("out of stock: %s", item)
		}
		inv.stock[item]--

	}
	return nil
}

func main() {
	inv := &Inventory{stock: map[string]int{"apple": 5, "banana": 3}}

	fmt.Println("order 1:", inv.Reserve([]string{"apple"}))
	fmt.Println("order 2:", inv.Reserve([]string{"apple", "banana"}))
	fmt.Println("all orders done")
}
