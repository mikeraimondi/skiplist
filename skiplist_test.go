package skiplist

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	maxLevel := 8
	list, err := New(maxLevel)
	if err != nil {
		t.Fatalf("expected no error calling New. got %s", err)
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
	// TODO insert nil,nil
	// TODO insert nil,val
	// TODO insert val,nil
	// TODO insert []byte{},[]byte{}
	// TODO insert []byte{},nil
	// TODO insert val,[]byte{}
	tests := []struct {
		name         string
		keyVals      [][]string
		expectedList *List
	}{
		{
			"with 1 pair",
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
			"with 2 sequential pairs",
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
		t.Run(tt.name, func(t *testing.T) {
			list := newList(2)

			for _, kv := range tt.keyVals {
				list.Insert([]byte(kv[0]), []byte(kv[1]))
			}

			compareNodes(t, tt.expectedList.header, list.header)
		})
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

	if bytes.Compare(expected.Value, actual.Value) != 0 {
		t.Fatalf("wrong node value. expected %q. got %q",
			expected.Value, actual.Value)
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

type testRandGen struct{}

func (s *testRandGen) Float32() float32 {
	return 0.9
}

func newList(maxLevel int) *List {
	header := &Node{
		Forward: make([]*Node, maxLevel, maxLevel),
	}

	return &List{
		maxLevel: maxLevel,
		header:   header,
		randGen:  &testRandGen{},
		less: func(a, b []byte) bool {
			return bytes.Compare(a, b) == -1
		},
	}
}
