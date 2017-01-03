package chring

import "testing"

func TestSearchFuncs(t *testing.T) {
	r := NewRing()
	r.Add("Bar")
	r.Add("Raz")
	r.Add("Qux")
	r.Add("Foo")

	/*
	   known hashIDs given default hasher
	   Bar -> 1320340042
	   Raz -> 1548738824
	   Qux -> 2661935400
	   Foo -> 3023971265
	*/

	tests := []struct {
		ID    string
		Index int
	}{
		{"Bar", 0},
		{"Raz", 1},
		{"Qux", 2},
		{"Foo", 3},
	}

	for _, test := range tests {
		if got, want := r.findNode(test.ID), test.Index; got != want {
			t.Errorf("findNode: got [%d] for %q, want [%d]", got, test.ID, want)
		}
	}

	for _, test := range tests {
		if got, want := r.searchByHashID(r.Hasher(test.ID)), test.Index; got != want {
			t.Errorf("searchByHashID: got [%d] for %q (#%d), want [%d]", got, test.ID, r.Hasher(test.ID), want)
		}
	}
}
