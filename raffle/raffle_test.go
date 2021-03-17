package raffle

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/hujw77/raffle/comparator"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

}

func testTreeMap(t *testing.T) {
	m := treemap.NewWith(comparator.BigIntComparator)
	m.Put(big.NewInt(1), "a")
	m.Put(big.NewInt(3), "b")
	m.Put(big.NewInt(6), "c")
	// (0, 1]  (1, 3] (3, 6]
	// k0, v0 := m.Ceiling(big.NewInt(0))
	_, v1 := m.Ceiling(big.NewInt(1))
	if actualValue, expectedValue := v1, 'a'; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	_, v2 := m.Ceiling(big.NewInt(2))
	if actualValue, expectedValue := v2, 'b'; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	_, v3 := m.Ceiling(big.NewInt(3))
	if actualValue, expectedValue := v3, 'b'; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	_, v4 := m.Ceiling(big.NewInt(4))
	if actualValue, expectedValue := v4, 'c'; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	_, v5 := m.Ceiling(big.NewInt(5))
	if actualValue, expectedValue := v5, 'c'; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	_, v6 := m.Ceiling(big.NewInt(6))
	if actualValue, expectedValue := v6, 'c'; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	// k7, v7 := m.Ceiling(big.NewInt(7))
}

func readCsvFile(filePath string) map[int]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	count := 0
	m := make(map[int]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		for value := range record {
			m[count] = record[value]
		}
		count++
	}

	return m
}

func U256(v string) *big.Int {
	v = strings.TrimPrefix(v, "0x")
	bn := new(big.Int)
	n, _ := bn.SetString(v, 16)
	return n
}

func TestRaffle(t *testing.T) {
	tickets := []*Ticket{
		NewTicket("0x1111111111111111111111111111111111111111", new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))),
		NewTicket("0x2222222222222222222222222222222222222222", new(big.Int).Mul(big.NewInt(2), big.NewInt(1e18))),
		NewTicket("0x3333333333333333333333333333333333333333", new(big.Int).Mul(big.NewInt(3), big.NewInt(1e18))),
		NewTicket("0x4444444444444444444444444444444444444444", new(big.Int).Mul(big.NewInt(4), big.NewInt(1e18))),
	}
	hashs := []*big.Int{
		U256("0x0000051d10559b5907127a70bc89d6ac20ad6204b865a862063a4b0d1095ae60"),
		U256("0x0000067182592e2885dd120d0b6212c254c861b9cd2b5c97bb594b36133eb54c"),
	}
	l, _ := NewLottery(tickets, hashs, 2)
	l.Pick()
	fmt.Println(l.Luckers())
}

// go test -v  -test.bench=".*"
func BenchmarkPick(b *testing.B) {
	tickets := []*Ticket{
		NewTicket("0x1111111111111111111111111111111111111111", new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))),
		NewTicket("0x2222222222222222222222222222222222222222", new(big.Int).Mul(big.NewInt(2), big.NewInt(1e18))),
		NewTicket("0x3333333333333333333333333333333333333333", new(big.Int).Mul(big.NewInt(3), big.NewInt(1e18))),
		NewTicket("0x4444444444444444444444444444444444444444", new(big.Int).Mul(big.NewInt(4), big.NewInt(1e18))),
	}
	mm := readCsvFile("../csv/query_result.csv")
	m := make(map[string]int)
	for i := 0; i < 500000; i++ {
		hashs := []*big.Int{
			U256(mm[i]),
			U256(mm[i+1]),
		}
		l, _ := NewLottery(tickets, hashs, 2)
		l.Pick()
		for _, addr := range l.Luckers() {
			m[addr]++
		}
	}
	fmt.Println(m)
}
