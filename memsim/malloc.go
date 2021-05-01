package memsim

import (
	"log"
)

func (h *Heap) checkAvail() {
	log.Printf("checkAvail() %v\n", h)
	slot := h.state[Slot]
	want := h.state[Want]
	idx := slotToIdx(slot)
	name := h.state[Name]
	h.resetState()

	if h.heads[idx] != NullPtr {
	// Available memory in this slot
		if slot == want {
		// Found a cell that we wanted; remove it
			h.state[Type] = SetHead
			h.state[Slot] = slot
			h.state[Loc] = h.heads[idx]
			h.state[Name] = name
		} else {
		// Need to split, since what we want is smaller
			h.state[Type] = Split
			h.state[Slot] = slot
			h.state[Want] = want
			h.state[Loc] = h.heads[idx]
			h.state[Name] = name
		}
	} else if slot >= MaxPwr {
	// Reached largest slot, but no head i.e. out of memory
		h.state[Type] = OutOfMem
	} else {
	// Currently nothing for this slot, so try to borrow from the next one
		h.state[Type] = CheckAvail
		h.state[Slot] = slot + 1
		h.state[Want] = want
		h.state[Name] = name
	}
}

func (h *Heap) split() {
	log.Printf("split() %v\n", h)
	slot := uint(h.state[Slot])
	want := uint(h.state[Want])
	name := h.state[Name]
	h.resetState()

	idx := slotToIdx(slot)
	loc := h.heads[idx]
	newSlot := slot - 1
	var shift uint = 1 << newSlot

	h.removeCell(loc)
	h.insertCell(loc+shift, newSlot)
	h.insertCell(loc, newSlot)
	if newSlot == want {
	// Split to desired slot
		h.state[Type] = SetHead
		h.state[Slot] = want
		h.state[Loc] = loc
		h.state[Name] = name
	} else {
	// Still need to split more
		h.state[Type] = Split
		h.state[Slot] = newSlot
		h.state[Want] = want
		h.state[Loc] = loc
		h.state[Name] = name
	}
}

func (h *Heap) setHead() {
	log.Printf("setHead() %v\n", h)
	loc := h.state[Loc]
	name := h.state[Name]
	slot := h.state[Slot]
	h.resetState()

	h.removeCell(loc)
	hdr := h.readHeader(loc)
	hdr.used = true
	hdr.slot = slot
	h.writeHeader(loc, hdr)

	h.vars[name] = loc + 1
	h.state[Type] = Idle
}
