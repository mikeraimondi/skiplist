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
	nilNode := newNil()
	tests := []struct {
		insertedKeys []uint
		insertedVals [][]byte
		expectedList *List
	}{
		{
			[]uint{1},
			[][]byte{[]byte("testing")},
			&List{
				Header: &Node{
					Forward: []*Node{
						&Node{
							Key:   uint(1),
							Value: []byte("testing"),
							Forward: []*Node{
								nilNode,
							},
						},
						nilNode,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		list := New(2)
		list.randGen = &testRandGen{}

		for i := 0; i < len(tt.insertedKeys); i++ {
			list.Insert(tt.insertedKeys[i], tt.insertedVals[i])
		}

		expectedNode := tt.expectedList.Header
		actualNode := list.Header
		compareNodes(t, expectedNode, actualNode)
	}
}

func compareNodes(t *testing.T, expected, actual *Node) {
	if expected == nil || actual == nil {
		t.Fatalf("nil node pointer!")
	}

	if expected.isNil() {
		if actual.isNil() {
			return
		}
		t.Fatalf("expected NIL node. got %+v", actual)
	}

	if expected.Key != actual.Key {
		t.Fatalf("wrong node key. expected %d. got %d",
			expected.Key, actual.Key)
	}

	expectedVal := expected.Value
	actualVal := actual.Value

	if len(expectedVal) != len(actualVal) {
		t.Fatalf("wrong node value. expected %q. got %q.",
			expectedVal,
			actualVal)
	}
	for i, element := range expectedVal {
		if element != actualVal[i] {
			t.Errorf("wrong node value at position %d. expected %q. got %q",
				i, element, actualVal[i])
		}
	}

	expectedForward := expected.Forward
	actualForward := actual.Forward

	if len(expectedForward) != len(actualForward) {
		t.Fatalf("wrong forward links. expected %v. got %v.",
			expectedForward,
			actualForward)
	}
	for i, node := range expectedForward {
		compareNodes(t, node, actualForward[i])
	}
}
