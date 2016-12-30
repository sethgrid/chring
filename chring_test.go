package chring_test

import (
	"log"
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

func TestAddIsConsistent(t *testing.T) {
	ring := chring.New()
	ring.Add("a")
	ring.Add("a")
	ring.Add("a")

	if ring.Nodes.Len() != 1 {
		t.Errorf("want 1 node, got %d", ring.Nodes.Len())
	}
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

func TestDrawChart(t *testing.T) {
	ring := chring.New()
	for _, n := range []string{"123.45.83.190", "123.45.83.191", "123.45.83.192", "123.45.78.191", "123.45.83.190", "123.45.83.190", "123.45.83.190", "123.45.83.190"} {
		ring.Add(n)
	}
	exampleURL := "http://localhost:8080/?key[]=%22user180%22&key[]=%22foo%22&hashid[]=5000000000000"
	log.Println(exampleURL)
	log.Println("blocking so we can manually test the url ^^")
	chring.Serve(ring)
}

func In(s string, slice []string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}
