// Package comparator provides ...
package comparator

import (
	"math/big"
)

func BigIntComparator(a, b interface{}) int {
	aAsserted := a.(*big.Int)
	bAsserted := b.(*big.Int)
	return aAsserted.Cmp(bAsserted)
}
