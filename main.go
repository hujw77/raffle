// Package main provides ...
package main

import (
	"fmt"
	"math/big"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/hujw77/raffle/comparator"
	"github.com/hujw77/raffle/raffle"
)

func testTreeMap() {
	m := treemap.NewWith(comparator.BigIntComparator)
	m.Put(big.NewInt(1), "x")
	m.Put(big.NewInt(3), "b")
	m.Put(big.NewInt(6), "a")
	// (0, 1]  (1, 3] (3, 6]
	// k0, v0 := m.Ceiling(big.NewInt(0))
	// fmt.Println(k0, v0)
	k1, v1 := m.Ceiling(big.NewInt(1))
	fmt.Println(k1, v1)
	k2, v2 := m.Ceiling(big.NewInt(2))
	fmt.Println(k2, v2)
	k3, v3 := m.Ceiling(big.NewInt(3))
	fmt.Println(k3, v3)
	k4, v4 := m.Ceiling(big.NewInt(4))
	fmt.Println(k4, v4)
	k5, v5 := m.Ceiling(big.NewInt(5))
	fmt.Println(k5, v5)
	k6, v6 := m.Ceiling(big.NewInt(6))
	fmt.Println(k6, v6)
	// k7, v7 := m.Ceiling(big.NewInt(7))
	// fmt.Println(k7, v7)
}

func main() {
	// testTreeMap()
	l, _ := raffle.New("http://localhost:7545", 100000)
	l.Print()
	l.Pick()
	l.Print()
	fmt.Println(l.Luckers())
}
