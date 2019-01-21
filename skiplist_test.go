package skiplist

import (
	"testing"
)

type testRandGen struct{}

func (s *testRandGen) Float32() float32 {
	return 0.9
}

func TestNew(t *testing.T) {
	maxLevel := 8
	list := New(maxLevel)

	if list.Level != 0 {
		t.Fatalf("wrong level for list. got %d. expected 0.", list.Level)
	}
	if list.MaxLevel != maxLevel {
		t.Fatalf("wrong maxLevel for list. got %d. expected %d.",
			list.MaxLevel,
			maxLevel)
	}

	header := list.Header

	if l := len(header.Forward); l != maxLevel {
		t.Fatalf("wrong length for header links. got %d. expected %d.", l, maxLevel)
	}

	nilNode := header.Forward[0]
	if nilNode.isNil() != true {
		t.Fatal("header not linked to NIL node")
	}
	for i := 0; i < maxLevel; i++ {
		if node := header.Forward[i]; node != nilNode {
			t.Fatalf("expected header link %d to be NIL node. got %+v", i, node)
		}
	}
}
func TestInsert(t *testing.T) {
	list := New(8)
	list.randGen = &testRandGen{}

	expectedKey := uint(1)
	expectedVal := []byte("testing")
	expectedNode := &Node{
		Key:   expectedKey,
		Value: expectedVal,
	}
	list.Insert(expectedKey, expectedVal)

	actualNode := list.Header.Forward[0]
	compareNodeContents(t, expectedNode, actualNode)
}

func compareNodeContents(t *testing.T, expected, actual *Node) {
	if expected.Key != actual.Key {
		t.Fatalf("wrong node key. expected %d. got %d",
			expected.Key, actual.Key)
	}

	expectedVal := expected.Value
	actualVal := actual.Value

	if len(expectedVal) != len(actualVal) {
		t.Fatalf("wrong node value length. expected %d. got %d.",
			len(expectedVal),
			len(actualVal))
	}
	for i, element := range expectedVal {
		if element != actualVal[i] {
			t.Errorf("wrong node value at position %d. expected %q. got %q",
				i, element, actualVal[i])
		}
	}
}
