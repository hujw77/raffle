// Package raffle provides ...
package raffle

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/hujw77/raffle/comparator"
)

type Ticket struct {
	user    common.Address
	balance *big.Int
}

type Lottery struct {
	tickets    *treemap.Map
	luckers    *treemap.Map
	quota      int64
	finalBlock uint64
	client     *ethclient.Client
}

func New(uri string, finalBlock uint64) (*Lottery, error) {
	client, err := ethclient.Dial(uri)
	if err != nil {
		return nil, err
	}
	list := []Ticket{
		Ticket{
			user:    common.HexToAddress("0x1111111111111111111111111111111111111111"),
			balance: new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)),
		},
		Ticket{
			user:    common.HexToAddress("0x2222222222222222222222222222222222222222"),
			balance: new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Ether)),
		},
		Ticket{
			user:    common.HexToAddress("0x3333333333333333333333333333333333333333"),
			balance: new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)),
		},
		Ticket{
			user:    common.HexToAddress("0x4444444444444444444444444444444444444444"),
			balance: new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Ether)),
		},
		Ticket{
			user:    common.HexToAddress("0x5555555555555555555555555555555555555555"),
			balance: new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether)),
		},
	}
	m := treemap.NewWith(comparator.AddressComparator)
	for _, ticket := range list {
		user, balance := ticket.user, ticket.balance
		if balance.Cmp(common.Big0) < 1 {
			return nil, fmt.Errorf("invalid user=%s, balance=%s", user, balance)
		}
		m.Put(user, balance)
	}
	return &Lottery{m, treemap.NewWith(comparator.AddressComparator), 1, finalBlock, client}, nil
}

func (l *Lottery) Size() int {
	return l.tickets.Size()
}

func (l *Lottery) convert() (*treemap.Map, error) {
	m := treemap.NewWith(comparator.BigIntComparator)
	l.tickets.Each(func(key interface{}, value interface{}) {
		user, balance := key.(common.Address), value.(*big.Int)
		last := common.Big0
		if m.Size() > 0 {
			max, _ := m.Max()
			last = max.(*big.Int)
		}
		m.Put(new(big.Int).Add(last, balance), user)
	})
	if l.Size() != m.Size() {
		return nil, errors.New("convert to treemap failed")
	}
	return m, nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (l *Lottery) Pick() {
	for i := int64(0); i < l.quota; i++ {
		l.pick()
	}
}

func (l *Lottery) random(max *big.Int) (*big.Int, error) {
	// block hash
	// block, err := l.client.BlockByNumber(context.Background(), new(big.Int).SetUint64(l.finalBlock))
	// if err != nil {
	// 	return nil, err
	// }
	// if block.NumberU64() != l.finalBlock {
	// 	return nil, errors.New("BlockByNumber returned wrong block")
	// }
	// hash := block.Hash()
	// fmt.Println("hash:", hash)
	// n := new(big.Int).Mod(hash.Big(), max)

	// math/rand
	rs := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	n := new(big.Int).Rand(rs, max)

	n = new(big.Int).Add(n, common.Big1)
	fmt.Println("n  :", n)
	return n, nil
}

func (l *Lottery) pick() (bool, error) {
	if l.Finished() {
		return true, nil
	}
	m, err := l.convert()
	if err != nil {
		return false, err
	}
	mx, _ := m.Max()
	max := mx.(*big.Int)
	fmt.Println("max:", max)
	d := new(big.Int).Div(max, big.NewInt(l.quota))
	fmt.Println("d  :", d)

	n, err := l.random(max)
	if err != nil {
		return false, err
	}

	k, v := m.Ceiling(n)
	fmt.Println("k  :", k, " v  :", v)
	l.luckers.Put(v, k)
	l.tickets.Remove(v)
	finished := l.Finished()
	if finished {
		return true, nil
	} else {
		l.incre()
		return false, nil
	}
	return l.Finished(), nil
}

func (l *Lottery) Print() {
	fmt.Println("luckers: ", l.luckers)
	fmt.Println("tickets: ", l.tickets)
	fmt.Println("quota: ", l.quota)
	fmt.Println("final: ", l.finalBlock)
}

func (l *Lottery) Finished() bool {
	return l.quota == int64(l.luckers.Size())
}

func (l *Lottery) Luckers() []common.Address {
	luckers := []common.Address{}
	for _, addr := range l.luckers.Keys() {
		luckers = append(luckers, addr.(common.Address))
	}
	return luckers
}

func (l *Lottery) incre() {
	l.finalBlock++
}
