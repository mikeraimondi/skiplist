package skiplist

import (
	"math/rand"
	"time"
)

const (
	maxKey = ^uint(0) // TODO should be greater than any legal key
	p      = 0.5
)

type randGen interface {
	Float32() float32
}

type Node struct {
	Key     uint
	Value   []byte
	Forward []*Node
}

func newNil() *Node {
	return &Node{
		Key: maxKey,
	}
}

func newNode(level int, key uint, value []byte) *Node {
	return &Node{
		Key:     key,
		Value:   value,
		Forward: make([]*Node, level+1, level+1),
	}
}

func (n *Node) isNil() bool {
	return n.Key == maxKey
}

type List struct {
	Header   *Node
	Level    int
	MaxLevel int
	randGen  randGen
}

func New(maxLevel int) *List {
	nilNode := newNil()
	header := &Node{
		Forward: make([]*Node, maxLevel, maxLevel),
	}
	for i := 0; i < maxLevel; i++ {
		header.Forward[i] = nilNode
	}

	randSrc := rand.NewSource(time.Now().Unix()) // TODO don't use timestamp
	return &List{
		Level:    0,
		MaxLevel: maxLevel,
		Header:   header,
		randGen:  rand.New(randSrc),
	}
}

func (l *List) Insert(searchKey uint, newValue []byte) {
	update := make([]*Node, l.MaxLevel)
	current := l.Header

	for i := l.Level; i >= 0; i-- {
		for current.Forward[i].Key < searchKey {
			current = current.Forward[i]
		}
		update[i] = current
	}

	current = current.Forward[0]
	if current.Key == searchKey {
		current.Value = newValue
		return
	}

	level := l.randomLevel()
	if level > l.Level {
		for i := l.Level + 1; i <= level; i++ {
			update[i] = l.Header
		}
		l.Level = level
	}
	node := newNode(level, searchKey, newValue)
	for i := 0; i <= level; i++ {
		node.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = node
	}
}

func (l *List) randomLevel() int {
	level := 0
	// TODO seed at list init time. allow tests to swap out random source
	for random := l.randGen.Float32(); random < p && level < l.MaxLevel; random = l.randGen.Float32() {
		level++
	}

	return level
}
