package memsim

import (
	"log"
)

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
			h.state[Loc] = h.heads[idx]
		} else {
			// Need to split
			h.state[Type] = Split
			h.state[Slot] = slot
			h.state[Want] = want
			h.state[Loc]  = h.heads[idx]
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
		h.state[Loc]  = loc
	}
}

func (h *Heap) setHead() {
	log.Printf("before setHead() %v\n", h)
	loc := h.state[Loc]
	h.resetState()

	h.removeCell(loc)
	hdr := h.readHeader(loc)
	hdr.used = true
	h.writeHeader(loc, hdr)
	h.state[Type] = Idle
	log.Printf("after setHead() %v\n", h)
}