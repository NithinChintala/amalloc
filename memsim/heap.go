package memsim

import (
	"log"
	"fmt"
	"os"
)

// Finite State Machine states
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
	MaxMalloc = (1 << MaxPwr) - HdrSize
)

// Specific type of state args
const (
	Type = "type"
	Slot = "slot"
	Want = "want"
	Loc  = "loc"
	Bdy  = "buddy"
	Name = "name"
)

type Heap struct {
	mem       []byte
	heads     []uint
	state     map[string]uint
	prevState map[string]uint
	vars      map[uint]uint
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
	// Insert in the memory
	h.insertCell(0, MaxPwr)

	h.state = make(map[string]uint)
	h.prevState = make(map[string]uint)
	h.vars = make(map[uint]uint)
	h.state[Type] = Idle
	h.prevState[Type] = Idle

	return &h
}

// Malloc allocates `size` amount of memory
func (h *Heap) Malloc(name string, size uint) {
	if size < 1 || size > MaxMalloc {
		fmt.Printf("Malloc(%d) is invalid, 0 < size <= %d\n", size, MaxMalloc)
		os.Exit(1)
	}
	log.Printf("Malloc(%d)\n", size)

	size += HdrSize
	slot := minPower2(size)
	char := strToUint(name)
	if _, ok := h.vars[char]; ok {
		fmt.Printf("Variable '%c' is already declared\n", char)
		os.Exit(1)
	}

	h.resetState()
	h.state[Type] = CheckAvail
	h.state[Slot] = slot
	h.state[Want] = slot
	h.state[Name] = strToUint(name)
}

func (h *Heap) Free(variable string) {
	p, ok := h.vars[strToUint(variable)]
	if !ok {
		fmt.Printf("Undeclared variabale given to Free(%s)\n", variable)
		os.Exit(1)
	}
	log.Printf("Free(%s) @ %d\n", variable, p)

	h.resetState()
	h.state[Type] = FreeSet
	h.state[Loc] = p - 1
	h.state[Slot] = h.readHeader(p - 1).slot
	h.state[Name] = strToUint(variable)
}

// Steps the malloc simulator one state
func (h *Heap) Step() {
	switch h.state[Type] {
	case Idle:
		h.idle()
	case Split:
		h.split()
	case CheckAvail:
		h.checkAvail()
	case SetHead:
		h.setHead()
	case ValSet:
	case FreeSet:
		h.freeSet()
	case BuddyChk:
		h.buddyChk()
	case BuddyFail:
		h.buddyFail()
	case BuddyMerge:
		h.buddyMerge()
	case OutOfMem:
		h.outOfMem()
	}
}

// Do nothing, stay idle
func (h *Heap) idle() {
	h.resetState()
	h.state[Type] = Idle
}

// Crash if out of memory
func (h *Heap) outOfMem() {
	fmt.Println("Heap is out of memory")
	os.Exit(1)
}
