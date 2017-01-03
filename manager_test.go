package chring_test

import (
	"testing"

	"github.com/sethgrid/chring"
)

func TestManagerAdd(t *testing.T) {
	ringManager := chring.NewRingManager()
	_ = ringManager.AddNode("node a")
	_ = ringManager.AddNode("node b")
	_ = ringManager.AddKey("user 180")
	_ = ringManager.AddKey("user 9")

	/*
		    The ringManager's data ring look like:
			[0] &{ID:user 180 HashID:870583203} (key)
			[1] &{ID:node b HashID:1413374556}  (node)
			[2] &{ID:user 9 HashID:3310596203}  (key)
			[3] &{ID:node a HashID:3442947046}  (node)
			This means that we should have one key in each node
	*/

	keysInA, _ := ringManager.GetKeys("node a")
	keysInB, _ := ringManager.GetKeys("node b")

	if len(keysInA) != 1 {
		for _, n := range keysInA {
			t.Logf("%#v", n)
		}
		t.Fatalf("got %d keys, want 1 key in node a", len(keysInA))
	}
	if got, want := keysInA[0].ID, "user 180"; got != want {
		t.Errorf("got %q, want %q as the key in node a", got, want)
	}
	if len(keysInB) != 1 {
		t.Fatalf("got %d keys, want 1 key in node b", len(keysInB))
	}
	if got, want := keysInB[0].ID, "user 9"; got != want {
		t.Errorf("got %q, want %q as the key in node b", got, want)
	}
}
