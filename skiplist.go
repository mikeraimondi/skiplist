package skiplist

import (
	"math/rand"
	"time"
)

const p = 0.5

type randGen interface {
	Float32() float32
}

type Node struct {
	Key     uint
	Value   []byte
	Forward []*Node
}

func newNode(level int, key uint, value []byte) *Node {
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
}

func New(maxLevel int) *List {
	header := &Node{
		Forward: make([]*Node, maxLevel, maxLevel),
	}

	randSrc := rand.NewSource(time.Now().Unix()) // TODO don't use timestamp
	return &List{
		level:    0,
		maxLevel: maxLevel,
		header:   header,
		randGen:  rand.New(randSrc),
	}
}

func (l *List) Insert(searchKey uint, newValue []byte) {
	update := make([]*Node, l.maxLevel)
	current := l.header

	for i := l.level; i >= 0; i-- {
		for l.less(current.Forward[i], searchKey) {
			current = current.Forward[i]
		}
		update[i] = current
	}

	current = current.Forward[0]
	if current != nil && current.Key == searchKey {
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

func (l *List) randomLevel() int {
	level := 0
	for random := l.randGen.Float32(); random < p && level < l.maxLevel; random = l.randGen.Float32() {
		level++
	}

	return level
}

func (l *List) less(node *Node, searchKey uint) bool {
	if node == nil {
		return false
	}
	return node.Key < searchKey
}
