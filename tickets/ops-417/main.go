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

	// create temp map/stock to hold new stock before commit to stock
	temp := make(map[string]int, len(items))

	for _, item := range items {
		avail, ok := temp[item]
		if !ok {
			avail = inv.stock[item]
		}

		if avail == 0 {
			return fmt.Errorf("out of stock: %s", item)
		}
		temp[item] = avail - 1
	}

	// once stock update is fixed without any error we can update the main stock in inventory
	for _, i := range items {
		inv.stock[i] = temp[i]
	}
	return nil
}

// 🎫 TICKET OPS-418 | Priority: High
// Reporter: Claude (PM)
// Assignee: Wazed

// Title: Stock numbers drifting, warehouse says counts are wrong

// Description:
// Since your fix shipped, nothing freezes anymore, great. But the warehouse team says some items show lower stock than they actually have.
// Customers got “out of stock” errors and their orders failed, yet stock still went down. Finance is now also yelling at me.

func main() {
	inv := &Inventory{stock: map[string]int{"apple": 3}}
	fmt.Println("order:", inv.Reserve([]string{"apple", "apple"}))
	fmt.Println("stock:", inv.stock)
}
