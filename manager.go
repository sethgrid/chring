package chring

import (
	"log"
	"sync"
)

/*
RingManager still a WIP
*/

type RingManager struct {
	sync.Mutex
	nodeNames  []string
	nodeRing   *Ring
	dataRing   *Ring
	keyFetcher func(nodeRing, dataRing *Ring, id string) (nodes, error)
	keyStorer  func(key string) error
}

func NewRingManager() *RingManager {
	r := NewRing()
	return &RingManager{
		nodeNames:  make([]string, 0),
		nodeRing:   NewRing(),
		dataRing:   r,
		keyFetcher: defaultKeyFetcher,
		keyStorer:  r.defaultKeyStorer,
	}
}

func (rm *RingManager) AddNode(nodeID string) error {
	rm.Lock()
	defer rm.Unlock()
	rm.nodeNames = append(rm.nodeNames, nodeID)
	rm.nodeRing.Add(nodeID)
	return rm.keyStorer(nodeID)
}

func (rm *RingManager) AddKey(key string) error {
	return rm.keyStorer(key)
}

func (rm *RingManager) GetKeys(nodeID string) (nodes, error) {
	return rm.keyFetcher(rm.nodeRing, rm.dataRing, nodeID)
}

// SetKeyFetcher allows a user to override the default in memory ring store
func (rm *RingManager) SetKeyFetcher(fn func(nodeRing, dataRing *Ring, id string) (nodes, error)) {
	rm.keyFetcher = fn
}

// SetKeyStorer allows a user to override the default in memory key store
func (rm *RingManager) SetKeyStorer(fn func(key string) error) {
	rm.keyStorer = fn
}

// Debug if true, prints verbose logging
var Debug = false

func debugf(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}

func defaultKeyFetcher(nodeRing, dataRing *Ring, id string) (nodes, error) {
	// r.Lock()
	// defer r.Unlock()

	startIndex := nodeRing.findNode(id)
	endIndex := startIndex + 1
	wraps := false
	if endIndex >= len(nodeRing.Nodes) {
		endIndex = 0
		wraps = true
	}
	debugf("looking for %q in", id)
	for i := 0; i < len(nodeRing.Nodes); i++ {
		debugf(">> node ring %+v", nodeRing.Nodes[i])
	}

	debugf("\nstartIndex (node %q): %d\nendIndex (the next node): %d", id, startIndex, endIndex)
	debugf("node ring length: %d", len(nodeRing.Nodes))
	debugf("data ring length: %d", len(dataRing.Nodes))

	debugf("end := dataRing.searchByHashID(nodeRing.Nodes[endIndex].HashID)")
	debugf("end := dataRing.searchByHashID(nodeRing.Nodes[%d].HashID)", endIndex)
	debugf("end := dataRing.searchByHashID(%d)", nodeRing.Nodes[endIndex].HashID)
	debugf("end := %d", dataRing.searchByHashID(nodeRing.Nodes[endIndex].HashID))

	start := dataRing.searchByHashID(nodeRing.Nodes[startIndex].HashID)
	end := dataRing.searchByHashID(nodeRing.Nodes[endIndex].HashID)

	debugf("parsing dataNodes. [%d] -> [%d]", start, end)
	for i := 0; i < len(dataRing.Nodes); i++ {
		debugf(">> data ring [%d] %+v", i, dataRing.Nodes[i])

	}

	// we subtract 2 because we don't count the start and end keys themselves as they are the node hashes, not key hashes
	size := start - end - 2
	if size < 0 {
		size = end - start - 2
	}
	dataNodes := make(nodes, size)

	if !wraps {
		debugf("does not wrap")
		for i := start + 1; i < end; i++ {
			debugf("appending [%d] %+v", i, dataRing.Nodes[i])
			dataNodes = append(dataNodes, dataRing.Nodes[i])
		}
	} else {
		debugf("wraps")
		for i := start + 1; i < len(dataRing.Nodes); i++ {
			debugf("appending [%d] %+v", i, dataRing.Nodes[i])
			dataNodes = append(dataNodes, dataRing.Nodes[i])
		}
		for i := 0; i < end; i++ {
			debugf("appending [%d] %+v", i, dataRing.Nodes[i])
			dataNodes = append(dataNodes, dataRing.Nodes[i])
		}
	}

	return dataNodes, nil
}

func (r *Ring) defaultKeyStorer(key string) error {
	r.Add(key)
	return nil
}
