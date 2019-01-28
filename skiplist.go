package skiplist

import (
	"bytes"
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
)

const p = 0.5

type randGen interface {
	Float32() float32
}

type Node struct {
	Key     []byte
	Value   []byte
	Forward []*Node
}

func newNode(level int, key, value []byte) *Node {
	return &Node{
		Key:     key,
		Value:   value,
		Forward: make([]*Node, level+1, level+1),
	}
}

type List struct {
	header   *Node
	level    int
	maxLevel int
	randGen  randGen
	less     func([]byte, []byte) bool
}

func New(maxLevel int) (*List, error) {
	header := &Node{
		Forward: make([]*Node, maxLevel, maxLevel),
	}

	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	randSrc := rand.NewSource(seed.Int64())

	return &List{
		maxLevel: maxLevel,
		header:   header,
		randGen:  rand.New(randSrc),
		less: func(a, b []byte) bool { // TODO allow callers to pass in `less`
			return bytes.Compare(a, b) == -1
		},
	}, nil
}

func (l *List) Insert(searchKey []byte, newValue []byte) {
	update := make([]*Node, l.maxLevel)
	current := l.header

	for i := l.level; i >= 0; i-- {
		for current.Forward[i] != nil && l.less(current.Forward[i].Key, searchKey) {
			current = current.Forward[i]
		}
		update[i] = current
	}

	current = current.Forward[0]
	if current != nil && (bytes.Compare(current.Key, searchKey) != 0) {
		current.Value = newValue
		return
	}

	level := l.randomLevel()
	if level > l.level {
		for i := l.level + 1; i <= level; i++ {
			update[i] = l.header
		}
		l.level = level
	}
	node := newNode(level, searchKey, newValue)
	for i := 0; i <= level; i++ {
		node.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = node
	}
}

func (l *List) Search(searchKey []byte) []byte {
	current := l.header

	for i := l.level; i >= 0; i-- {
		for current.Forward[i] != nil && l.less(current.Forward[i].Key, searchKey) {
			current = current.Forward[i]
		}
	}

	current = current.Forward[0]
	if current != nil && (bytes.Compare(current.Key, searchKey) == 0) {
		return current.Value
	}

	return []byte{}
}

func (l *List) randomLevel() int {
	level := 0
	for random := l.randGen.Float32(); random < p && level < l.maxLevel; random = l.randGen.Float32() {
		level++
	}

	return level
}
