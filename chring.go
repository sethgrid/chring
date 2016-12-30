package chring

import (
	"errors"
	"hash/crc32"
	"sort"
	"sync"
)

// Ring is a consistent hash ring. Use New() to create a ring. You may change out the hasher function to change key balancing if needed.
type Ring struct {
	sync.Mutex
	Nodes  nodes
	Hasher func(id string) uint32
}

// New creates a new consistent hash ring with a default hashing algo
func New() *Ring {
	return &Ring{Nodes: []*node{}, Hasher: DefaultHasher}
}

// Add inserts a new node into the hash ring
func (r *Ring) Add(id string) {
	r.Lock()
	defer r.Unlock()

	// don't insert the same node more than once
	if r.findNode(id) == 0 && len(r.Nodes) > 0 && r.Nodes[0].ID == id {
		// TODO: why is search returning 0 on miss? should be -1, yeah?
		return
	}

	n := newNode(id, r.Hasher)
	r.Nodes = append(r.Nodes, n)
	sort.Sort(r.Nodes)
}

// Get retrievs the closest node in the hash ring for the given key
func (r *Ring) Get(key string) string {
	r.Lock()
	r.Unlock()

	if len(r.Nodes) == 0 {
		return "" // should error?
	}

	i := r.search(key)
	if i >= r.Nodes.Len() || i == -1 {
		i = 0 // default to initial node
	}
	return r.Nodes[i].ID
}

var ErrNotFound = errors.New("node not found")

func (r *Ring) Remove(id string) error {
	r.Lock()
	defer r.Unlock()

	i := r.findNode(id)
	if i >= r.Nodes.Len() || i == -1 || r.Nodes[i].ID != id {
		return ErrNotFound
	}

	r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)
	return nil
}

// findNode is different than `search` in that it searches for an exact match for a node ID
func (r *Ring) findNode(id string) int {
	return sort.Search(len(r.Nodes), func(i int) bool {
		return r.Nodes[i].HashID == r.Hasher(id)
	})
}

// search is different than `findNode` in that it searches for any node next in the hash ring for a given key
func (r *Ring) search(key string) int {
	return sort.Search(len(r.Nodes), func(i int) bool {
		return r.Nodes[i].HashID < r.Hasher(key)
	})
}

// node comprises nodes, which are placed in the consistent hash ring
type node struct {
	ID     string
	HashID uint32
}

// newNode creates a new node to go into the hash ring
func newNode(id string, fn hasher) *node {
	return &node{
		ID:     id,
		HashID: fn(id),
	}
}

// DefaultHasher uses crc32
func DefaultHasher(id string) uint32 {
	return crc32.ChecksumIEEE([]byte(id))
}

// hasher is an alias type used in the Ring type
type hasher func(id string) uint32

// nodes is an alias type for easy reference for matching the swap interface
type nodes []*node

// Len() is for matching the swap interface
func (n nodes) Len() int { return len(n) }

// Swap() is for matching the swap interface
func (n nodes) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// Less() is for matching the swap interface
func (n nodes) Less(i, j int) bool { return n[i].HashID < n[j].HashID }
