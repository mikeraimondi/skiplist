package skiplist

import (
	"bytes"
	"testing"
)

type testRandGen struct{}

func (s *testRandGen) Float32() float32 {
	return 0.9
}

func TestNew(t *testing.T) {
	maxLevel := 8
	list, err := New(maxLevel)
	if err != nil {
		t.Fatalf("New returned error: %s", err)
	}

	if list.level != 0 {
		t.Fatalf("wrong level for list. got %d. expected 0.", list.level)
	}
	if list.maxLevel != maxLevel {
		t.Fatalf("wrong maxLevel for list. got %d. expected %d.",
			list.maxLevel,
			maxLevel)
	}

	header := list.header

	if l := len(header.Forward); l != maxLevel {
		t.Fatalf("wrong length for header links. got %d. expected %d.", l, maxLevel)
	}

	for i := 0; i < maxLevel; i++ {
		if node := header.Forward[i]; node != nil {
			t.Fatalf("expected header link %d to be nil. got %+v", i, node)
		}
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		keyVals      [][]string
		expectedList *List
	}{
		{
			[][]string{
				[]string{"foo", "bar"},
			},
			&List{
				header: &Node{
					Forward: []*Node{
						&Node{
							Key:   []byte("foo"),
							Value: []byte("bar"),
							Forward: []*Node{
								nil,
							},
						},
						nil,
					},
				},
			},
		},
		{
			[][]string{
				[]string{"a", "b"},
				[]string{"c", "d"},
			},
			&List{
				header: &Node{
					Forward: []*Node{
						&Node{
							Key:   []byte("a"),
							Value: []byte("b"),
							Forward: []*Node{
								&Node{
									Key:   []byte("c"),
									Value: []byte("d"),
									Forward: []*Node{
										nil,
									},
								},
							},
						},
						nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		list, err := New(2)
		if err != nil {
			t.Fatalf("New returned error: %s", err)
		}
		list.randGen = &testRandGen{}

		for _, kv := range tt.keyVals {
			list.Insert([]byte(kv[0]), []byte(kv[1]))
		}

		expectedNode := tt.expectedList.header
		actualNode := list.header
		compareNodes(t, expectedNode, actualNode)
	}
}

func compareNodes(t *testing.T, expected, actual *Node) {
	if expected == nil {
		if actual == nil {
			return
		}
		t.Fatalf("expected NIL node. got %+v", actual)
	}

	if bytes.Compare(expected.Key, actual.Key) != 0 {
		t.Fatalf("wrong node key. expected %q. got %q",
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
