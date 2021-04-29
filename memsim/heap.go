package memsim

import (
	"log"
)

// FSM States
const (
	Idle = iota
	Split
	CheckAvail
	SetHead
	ValSet
	FreeSet
	BuddyChk
	BuddyFail
	BuddyMerge
	OutOfMem
)

const (
	NullPtr = 0b00011111 // 31

	HdrSize  = 1
	CellSize = 2
	MinPwr   = 1
	MaxPwr   = 4
	NumSlots = MaxPwr - MinPwr + 1
)

// State args
const (
	Type = "type"
	Slot = "slot"
	Want = "want"
	Loc = "loc"
	Bdy = "buddy"
)

type Heap struct {
	mem   []byte
	heads []uint
	state map[string]uint
	prevState map[string]uint
}

type Pointer uint8

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
	h.heads[NumSlots-1] = 0
	// TODO write a insertCell func
	h.mem[0] = 0b01111111
	h.mem[1] = NullPtr
	//h.insertCell(0, MaxPwr)

	h.state = make(map[string]uint)
	h.prevState = make(map[string]uint)
	h.state[Type] = Idle
	h.prevState[Type] = Idle

	return &h
}

// Steps the malloc simulator one state
func (h *Heap) Step() {
	stateType := h.state[Type]
	switch stateType {
	case Idle:
		// Do nothing if Idle
		h.resetState()
		h.state[Type] = Idle
		return
	case Split:
		h.split()
		return
	case CheckAvail:
		h.checkAvail()
		return
	case SetHead:
		h.setHead()
		return
	case ValSet:
		return
	case FreeSet:
		h.freeSet()
		return
	case BuddyChk:
		h.buddyChk()
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
	if size < 1 || size >= maxMalloc {
		log.Fatalf("Malloc(%d) is invalid, 0 < size < %d", size, maxMalloc)
	}
	log.Printf("Malloc(%d)\n", size)

	size += HdrSize
	slot := minPower2(size)

	h.resetState()
	h.state[Type] = CheckAvail
	h.state[Slot] = slot
	h.state[Want] = slot
}

func (h *Heap) Free(p uint) {
	log.Printf("Free(%d)\n", p)

	h.resetState()
	h.state[Type] = FreeSet
	h.state[Loc] = p - 1
}
