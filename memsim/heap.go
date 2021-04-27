package memsim

import (
	"log"
	"fmt"
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
)

type Heap struct {
	mem   []byte
	heads []uint
	state map[string]uint
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

func (h *Heap) checkAvail() {
	slot := h.state[Slot]
	want := h.state[Want]
	idx := slotToIdx(slot)
	log.Printf("checkAvail() %v\n", h)
	h.resetState()
	if h.heads[idx] != NullPtr {
		if slot == want {
			// Found a cell that we wanted; remove it
			h.state[Type] = SetHead
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

	idx := slotToIdx(slot)
	loc := h.heads[idx]
	newSlot := slot - 1
	var shift uint = 1 << newSlot

	h.removeCell(loc)
	// TODO when you insert have to figure something out to not make these merge
	// immediately, probably just set the first cell to be "used"
	// Actually, might just want to manually call a merge() function instead
	// This feels really bad
	h.insertCell(loc+shift, newSlot, false)
	h.insertCell(loc, newSlot, true)
	if newSlot == want {
		// Split to desired slot
		h.state[Type] = SetHead
		h.state[Slot] = want
		h.state[Loc] = loc
	} else {
		// Still need to split more
		h.state[Type] = Split
		h.state[Slot] = newSlot
		h.state[Want] = want
	}
}

func (h *Heap) setHead() {
	log.Printf("setHead() %v\n", h)
	loc := h.state[Loc]
	h.resetState()

	h.removeCell(loc)
	hdr := h.readHeader(loc)
	hdr.used = true
	h.writeHeader(loc, hdr)
	h.state[Type] = Idle
}

// Update the pointers in the cell free list
func (h *Heap) removeCell(loc uint) {
	log.Printf("before removeCell(%d) %v\n", loc, h)
	oldCell := h.readCell(loc)
	idx := slotToIdx(oldCell.slot)

	if oldCell.prev == NullPtr {
	// Front of free list
		if oldCell.next == NullPtr {
		// singelton list
			log.Println("singleton list")
			h.heads[idx] = NullPtr
		} else {
		// front of list len > 1
			log.Println("front, len > 1")
			//fmt.Println(oldCell)
			newFront := h.readCell(oldCell.next)
			newFront.prev = NullPtr
			h.writeCell(oldCell.next, newFront)
			h.heads[idx] = oldCell.next
		}
	} else if oldCell.next == NullPtr {
	// At end of list len > 1
		log.Println("end, len > 1")
		prevCell := h.readCell(oldCell.prev)
		prevCell.next = oldCell.next
		h.writeCell(oldCell.prev, prevCell)
	} else {
	/// At Middle of list
		log.Println("middle")
		prevCell := h.readCell(oldCell.prev)
		nextCell := h.readCell(oldCell.next)

		prevCell.next = oldCell.next
		nextCell.prev = oldCell.prev

		h.writeCell(oldCell.prev, prevCell)
		h.writeCell(nextCell.prev, nextCell)
	}
	// Do this?
	// oldCell.used = true
	//h.writeCell(loc, oldCell)
	log.Printf("after removeCell(%d) %v\n", loc, h)
}

func (h *Heap) insertCell(loc, slot uint, used bool) {
	log.Printf("before insertCell(loc=%d, slot=%d, used=%t) %v\n", loc, slot, used, h)
	idx := slotToIdx(slot)
	newCell := Cell{}
	newCell.slot = slot
	newCell.used = used

	/*
	buddyLoc := h.getBuddy(loc)
	// Found a buddy, merge them + recursively insert
	if buddyLoc != NullPtr {
	}
	*/
	if oldFrontLoc := h.heads[idx]; oldFrontLoc != NullPtr {
	// The slot has something
		oldFront := h.readCell(oldFrontLoc)
		oldFront.prev = loc

		newCell.prev = NullPtr
		newCell.next = oldFrontLoc

		h.writeCell(oldFrontLoc, oldFront)
	} else {
	// The slot is empty
		newCell.prev = NullPtr
		newCell.next = NullPtr
	}
	h.writeCell(loc, newCell)
	h.heads[idx] = loc
	log.Printf("after insertCell(loc=%d, slot=%d, used=%t) %v\n", loc, slot, used, h)
}

func (h *Heap) insertCellMerge() {

}

func (h *Heap) Free(p int) {

}

func (h *Heap) resetState() {
	for k := range h.state {
		delete(h.state, k)
	}
}

func (h *Heap) getBuddy(loc uint) uint {
	hdr := h.readHeader(loc)
	if hdr.slot >= MaxPwr {
		// MaxPwr slot has not buddies
		return NullPtr
	}
	buddyLoc := loc ^ (1 << hdr.slot)
	buddy := h.readHeader(buddyLoc)
	if !buddy.used && buddy.slot == hdr.slot {
		// Found the buddy, has same size + is not used
		return buddyLoc
	}
	// Buddy is currently being used
	return NullPtr
}

func (h *Heap) logJustify(format string, a ...interface{}) {
	ctx := fmt.Sprintf("%-20s", fmt.Sprintf(format, a...))
	log.Printf("%s %v\n", ctx, h)
}
