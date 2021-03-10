// Package comparator provides ...
package comparator

import (
	"math/big"

	"github.com/emirpasic/gods/utils"
	"github.com/ethereum/go-ethereum/common"
)

func BigIntComparator(a, b interface{}) int {
	aAsserted := a.(*big.Int)
	bAsserted := b.(*big.Int)
	return aAsserted.Cmp(bAsserted)
}

func AddressComparator(a, b interface{}) int {
	aAsserted := a.(common.Address).String()
	bAsserted := b.(common.Address).String()
	return utils.StringComparator(aAsserted, bAsserted)
}
