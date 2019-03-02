package skiplist_test

import (
	"testing"

	"github.com/mikeraimondi/skiplist"
)

func TestInsertAndSearch(t *testing.T) {
	tests := []struct {
		name         string
		inserts      [][]string
		searches     []string
		expectedVals []string
		expectedOks  []bool
	}{
		{
			"with 1 pair, key found",
			[][]string{
				[]string{"foo", "bar"},
			},
			[]string{"foo"},
			[]string{"bar"},
			[]bool{true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := test
			t.Parallel()

			list, err := skiplist.New(2)
			if err != nil {
				t.Fatalf("expected no error initializing list. got %s", err)
			}

			for _, kv := range tt.inserts {
				list.Insert([]byte(kv[0]), []byte(kv[1]))
			}

			for i, search := range tt.searches {
				actualVal, actualOk := list.Search([]byte(search))

				if tt.expectedVals[i] != string(actualVal) {
					t.Fatalf("wrong value from Search(%q). expected %q. got %q.",
						search, tt.expectedVals[i], actualVal)
				}
				if tt.expectedOks[i] != actualOk {
					t.Fatalf("wrong OK from Search(%q). expected %t. got %t.",
						search, tt.expectedOks[i], actualOk)
				}
			}
		})
	}
}
