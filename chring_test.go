package chring_test

import (
	"testing"

	"github.com/sethgrid/chring"
)

var NodeList = []string{"node 1", "node 2", "node 3", "node 4"}

func newSeededRing() *chring.Ring {
	ring := chring.New()
	for _, n := range NodeList {
		ring.Add(n)
	}
	return ring
}

func TestGetReturnsAKey(t *testing.T) {
	ring := newSeededRing()

	for _, key := range []string{"user A", "user B", "user C", "user D"} {
		nodeName := ring.Get(key)
		if !In(nodeName, NodeList) {
			t.Errorf("got %q for Get(%q)", nodeName, key)
		}
	}
}

func TestDeletedNodesAreNotReturned(t *testing.T) {
	ring := newSeededRing()
	err := ring.Remove("node 1")
	if err != nil {
		t.Errorf("got error %q, want nil when removing a known node", err.Error())
	}

	for _, key := range []string{"user A", "user B", "user C", "user D"} {
		nodeName := ring.Get(key)
		if nodeName == "node 1" {
			t.Errorf("got %q for Get(%q), but it has been deleted", "node 1", key)
		}
	}
}

func TestDeleteUnknownNode(t *testing.T) {
	ring := newSeededRing()
	err := ring.Remove("node x") // known to not exist
	if err == nil {
		t.Error("got no error, but should when deleting a non-existant node")
	}
}

func TestConsistentHashRing(t *testing.T) {
	ring := newSeededRing()

	want := ring.Get("foo")
	errCount := 0
	for i := 0; i < 100; i++ {
		got := ring.Get("foo")
		if got != want {
			errCount++
		}
	}
	if errCount != 0 {
		t.Errorf("got %d inconsistent hash results, want 0 if we are consistent", errCount)
	}
}

func In(s string, slice []string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}
