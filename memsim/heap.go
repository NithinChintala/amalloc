package memsim

import (
	"log"
)

// FSM States
const (
	Idle = iota
	Split
	CheckAvail
	RemoveCell
	ValSet
	FreeSet
	BuddyChk
	BuddyFail
	BuddyMerge
	OutOfMem
)

const (
	// Why 255? Since I store the next pointer as on offset [0, 1 <<4]
	// in a uint8, 255 was the next best after 0, 255 = 0b11111111
	NullPtr = 255

	HdrSize = 1
	MinPwr  = 1
	MaxPwr  = 4
	NumSlots = MaxPwr - MinPwr + 1
)

// State args
const (
	Type = "type"
	Slot = "slot"
	Want = "want"
)

type Heap struct {
	mem     []byte
	heads   []uint
	state   map[string]int
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

func NewHeap() *Heap {
	h := Heap{}
	h.mem = make([]byte, 1<<MaxPwr)

	// Heads are orders from lowest to highest level
	// heads[i] would have allocation size of 1 << (MinPwr + i)
	h.heads = make([]uint, NumSlots)
	for i := range h.heads {
		h.heads[i] = NullPtr
	}
	// Largest level's head always starts at the beginning
	h.heads[NumSlots - 1] = 0
	// TODO write a insertCell func
	h.mem[0] = 0b10000011
	h.mem[1] = NullPtr

	h.state = make(map[string]int)
	h.state[Type] = Idle

	return &h
}

// Steps the malloc simulator one state
func (h *Heap) Step() {
	stateType := h.state[Type]
	switch stateType {
	case Idle:
		// Do nothing if Idle
		return
	case Split:
		return
	case CheckAvail:
		h.checkAvail()
		return
	case RemoveCell:
		return
	case ValSet:
		return
	case FreeSet:
		return
	case BuddyChk:
		return
	case BuddyFail:
		return
	case BuddyMerge:
		return
	case OutOfMem:
		return
	}
}

// Malloc allocates `size` amount of memory
// Returns an error if heap is out of memory
func (h *Heap) Malloc(size uint) {
	var maxMalloc uint = 1 << MaxPwr
	if size < 1 || size > maxMalloc {
		log.Fatalf("Malloc(%d) is invalid, 0 < size <= %d", size, maxMalloc)
	}
	log.Printf("Malloc(%d)\n", size)

	size += HdrSize
	slot := minPower2(size)

	h.resetState()
	h.state[Type] = CheckAvail
	h.state[Slot] = int(slot)
	h.state[Want] = int(slot)
}

func (h *Heap) checkAvail() {
	slot := h.state[Slot]
	want := h.state[Want]
	idx := slotToIdx(uint(slot))
	log.Printf("checkAvail() %v\n", h)
	h.resetState()
	if h.heads[idx] != NullPtr {
		if slot == want {
			// Found a cell that we wanted; remove it
			h.state[Type] = RemoveCell
			h.state[Slot] = slot
		} else {
			// Need to split
			h.state[Type] = Split
			h.state[Slot] = slot
			h.state[Want] = want
		}
	} else if slot >= MaxPwr {
		// Reached largest slot, but no head i.e. out of memory
		h.state[Type] = OutOfMem
	} else {
		// Currently nothing for this slot, so try to borrow from the next one
		h.state[Type] = CheckAvail
		h.state[Slot] = slot + 1
		h.state[Want] = want
	}
}

func (h *Heap) split() {
	slot := uint(h.state[Slot])
	want := uint(h.state[Want])
	log.Printf("split() %v\n", h)
	h.resetState()

	idx := slotToIdx(uint(slot))
	loc := h.heads[idx]
	newSlot := slot - 1
	var shift uint = 1 << newSlot

	h.removeCell(slot)
	h.insertCell(loc + shift, newSlot)

	if newSlot == want {
		// Split to desired slot
		h.state[Type] = Idle
	} else {
		// Still need to split more
		h.state[Type] = Split
		h.state[Slot] = int(newSlot)
		h.state[Want] = int(want)
	}
}

// Should always be removing from head???
func (h *Heap) removeCell(slot uint) {
}

func (h *Heap) insertCell(loc, slot uint) {

}


func (h *Heap) Free(pointer int) {

}

func (h *Heap) resetState() {
	for k := range h.state {
		delete(h.state, k)
	}
}

func minPower2(x uint) uint {
	var power uint = 0
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
