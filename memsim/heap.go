package memsim

import (
	"log"
)

// FSM States
const (
	Idle  = iota
	Split
	CheckAvail
	HeadSet
	ValSet
	FreeSet
	BuddyChk
	BuddyFail
	BuddyMerge
	OutOfMem
)

const (
	NullPtr = -1
	HdrSize = 1
	MinPwr = 1
)

// State args
const (
	Type = "type"
	Slot = "slot"
	Want = "want"
)

type Heap struct {
	// 1 << size, i.e. size=4 => 16 Bytes
	maxSlot uint
	mem []byte
	heads []int
	state map[string]int
}

type Cell struct {
	used bool
	slot uint
	next uint
}

type Header struct {
	used bool
	slot uint
}

func NewHeap(size uint) *Heap {
	if size <= HdrSize {
		log.Fatalf("NewHeap(%d) is invalid, size > %d\n", size, HdrSize)
	}
	minPwr := minPower2(size)

	h := Heap{}
	h.maxSlot = minPwr
	h.mem = make([]byte, 1 << minPwr)

	// Heads are orders from lowest to highest level
	// heads[i] would have allocation size of 1 << (MinPwr + i)
	h.heads = make([]int, slotToIdx(minPwr) + 1)
	for i := range h.heads {
		h.heads[i] = NullPtr
	}
	// Largest level's head always starts at the beginning
	h.heads[size - HdrSize] = 0

	h.state = make(map[string]int)
	h.state[Type] = Idle

	return &h
}

// Malloc allocates `size` amount of memory
// Returns an error if heap is out of memory
func (h *Heap) Malloc(size uint) {
	var maxMalloc uint = (1 << h.maxSlot) - HdrSize
	if size < 1 || size > maxMalloc {
		log.Fatalf("Malloc(%d) is invalid, 0 < size <= %d", size, maxMalloc)
	}
	size += HdrSize
	slot := minPower2(size)
	h.resetState()
	h.state[Type] = CheckAvail
	h.state[Slot] = int(slot)
	h.state[Want] = int(slot)

	log.Printf("Malloc(%d)\n", size)
}

func (h *Heap) checkAvail() {
	slot := uint(h.state[Slot])
	want := h.state[Want]
	idx := slotToIdx(slot)
	h.resetState()
	if h.heads[idx] != NullPtr {
		// Found a head so move on to the next state
		h.state[Type] = HeadSet
	} else if slot == h.maxSlot {
		// Reached largest slot, but no head i.e. out of memory
		h.state[Type] = OutOfMem
	} else {
		// Currently nothing for this slot, so try to borrow from the next one
		h.state[Type] = CheckAvail
		h.state[Slot] = int(slot + 1)
		h.state[Want] = want
	}
}

func (h *Heap) Free(pointer int) {

}

func (h *Heap) resetState() {
	for k := range h.state {
		delete(h.state, k)
	}
}

func minPower2(x uint) uint {
	var power uint = 1
	var curr uint = 1
	for curr <= x {
		power++
		curr <<= 1
	}
	return power
}

func slotToIdx(slot uint) uint {
	return slot - MinPwr
}
