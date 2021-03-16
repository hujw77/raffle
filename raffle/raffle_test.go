package raffle

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// go test -v  -test.bench=".*"
func BenchmarkPick(b *testing.B) {
	m := make(map[common.Address]int)
	blockN := uint64(2)
	for i := 0; i < 500000; i++ {
		l, _ := New("http://localhost:7545", blockN)
		l.Pick()
		for _, addr := range l.Luckers() {
			m[addr]++
		}
		blockN++
	}
	fmt.Println(m)
}
