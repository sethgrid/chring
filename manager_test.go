package chring_test

import (
	"testing"

	"github.com/sethgrid/chring"
)

func TestManager(t *testing.T) {
	ringManager := chring.NewRingManager()
	_ = ringManager.AddNode("node a")
	_ = ringManager.AddNode("node b")
	_ = ringManager.AddKey("user 180")
	_ = ringManager.AddKey("user 9")

	nodeNames := ringManager.GetNodes()
	if len(nodeNames) != 2 {
		for _, n := range nodeNames {
			t.Logf("%s", n)
		}
		t.Errorf("got %d, want 2 for count of node names", len(nodeNames))
	}

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

	/*
		Now remove the node b, we should have the following:
			[0] &{ID:user 180 HashID:870583203} (key)
			[1] &{ID:user 9 HashID:3310596203}  (key)
			[2] &{ID:node a HashID:3442947046}  (node)
	*/

	_ = ringManager.RemoveNode("node b")
	keysInA, _ = ringManager.GetKeys("node a")
	keysInB, err := ringManager.GetKeys("node b")

	if err != chring.ErrNotFound {
		t.Errorf("got error %v, want %v", err, chring.ErrNotFound)
	}

	if len(keysInA) != 2 {
		for _, n := range keysInA {
			t.Logf("%#v", n)
		}
		t.Fatalf("got %d keys, want 2 keys in node a", len(keysInA))
	}
	if got, want := keysInA[0].ID, "user 180"; got != want {
		t.Errorf("got %q, want %q as the key in node a", got, want)
	}
	if got, want := keysInA[1].ID, "user 9"; got != want {
		t.Errorf("got %q, want %q as the key in node a", got, want)
	}
	if len(keysInB) != 0 {
		for _, n := range keysInB {
			t.Logf("%#v", n)
		}
		t.Fatalf("got %d keys, want 0 keys in node b", len(keysInB))
	}

	/*
		Now remove both keys. We should have the following:
			[0] &{ID:node a HashID:3442947046}  (node)
	*/

	_ = ringManager.RemoveKey("user 180")
	_ = ringManager.RemoveKey("user 9")

	keysInA, _ = ringManager.GetKeys("node a")
	if len(keysInA) != 0 {
		for _, n := range keysInA {
			t.Logf("%#v", n)
		}
		t.Fatalf("got %d keys, want 0 keys in node a", len(keysInA))
	}
}
