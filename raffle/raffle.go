// Package raffle provides ...
package raffle

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/hujw77/raffle/comparator"
)

var Big0 = big.NewInt(0)
var Big1 = big.NewInt(1)

type Ticket struct {
	user    string
	balance *big.Int
}

type Lottery struct {
	tickets *treemap.Map
	luckers *treemap.Map
	hashs   []*big.Int
	quota   int
}

func NewTicket(user string, balance *big.Int) *Ticket {
	return &Ticket{user, balance}
}

func NewLottery(tickets []*Ticket, hashs []*big.Int, quota int) (*Lottery, error) {
	if len(hashs) != quota {
		return nil, errors.New("Invalid quota or hashs length")
	}
	m := treemap.NewWith(utils.StringComparator)
	for _, ticket := range tickets {
		user, balance := ticket.user, ticket.balance
		if balance.Cmp(Big0) < 1 {
			return nil, fmt.Errorf("invalid user=%s, balance=%s", user, balance)
		}
		m.Put(user, balance)
	}
	return &Lottery{m, treemap.NewWith(utils.StringComparator), hashs, quota}, nil
}

func (l *Lottery) Size() int {
	return l.tickets.Size()
}

func (l *Lottery) convert() (*treemap.Map, error) {
	m := treemap.NewWith(comparator.BigIntComparator)
	l.tickets.Each(func(key interface{}, value interface{}) {
		user, balance := key.(string), value.(*big.Int)
		last := Big0
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

func (l *Lottery) Pick() (bool, error) {
	for i := 0; i < l.quota; i++ {
		if err := l.pick(i); err != nil {
			return false, err
		}
	}
	return true, nil
}

// func init() {
// 	rand.Seed(time.Now().UTC().UnixNano())
// }

func (l *Lottery) random(max *big.Int, i int) *big.Int {
	hash := l.hashs[i]
	n := new(big.Int).Mod(hash, max)

	// math/rand
	// rs := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// n := new(big.Int).Rand(rs, max)

	n = new(big.Int).Add(n, Big1)
	return n
}

func (l *Lottery) pick(i int) error {
	m, err := l.convert()
	if err != nil {
		return err
	}
	mx, _ := m.Max()
	max := mx.(*big.Int)
	n := l.random(max, i)
	k, v := m.Ceiling(n)
	l.luckers.Put(v, k)
	l.tickets.Remove(v)
	return nil
}

func (l *Lottery) Print() {
	fmt.Println("luckers: ", l.luckers)
	fmt.Println("tickets: ", l.tickets)
	fmt.Println("hashs: ", l.hashs)
	fmt.Println("quota: ", l.quota)
}

func (l *Lottery) Finished() bool {
	return l.quota == l.luckers.Size()
}

func (l *Lottery) Luckers() []string {
	luckers := []string{}
	for _, addr := range l.luckers.Keys() {
		luckers = append(luckers, addr.(string))
	}
	return luckers
}
