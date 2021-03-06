package memsim

import (
	"log"
)

func (h *Heap) freeSet() {
	log.Printf("freeSet() %v\n", h)
	loc := h.state[Loc]
	name := h.state[Name]
	h.resetState()

	hdr := h.readHeader(loc)
	h.insertCell(loc, hdr.slot)
	delete(h.vars, name)

	bdy := h.getBuddy(loc)
	h.state[Type] = BuddyChk
	h.state[Loc] = loc
	h.state[Bdy] = bdy
}

func (h *Heap) buddyChk() {
	log.Printf("buddyChk() %v\n", h)
	loc := h.state[Loc]
	bdy := h.state[Bdy]
	h.resetState()

	// Called free() on slot 4 variable
	hdr := h.readHeader(loc)
	if bdy == NullPtr {
		h.state[Type] = BuddyFail
		h.state[Loc] = loc
		h.state[Bdy] = loc
		h.state[Slot] = hdr.slot
		return
	}

	buddy := h.readHeader(bdy)
	if !buddy.used && buddy.slot == hdr.slot {
	// Found the buddy, has same size + is not used
		h.state[Type] = BuddyMerge
		h.state[Loc] = loc
		h.state[Slot] = hdr.slot
		h.state[Bdy] = bdy
	} else {
	// Buddy is not free
		h.state[Type] = BuddyFail
		h.state[Loc] = loc
		h.state[Bdy] = bdy
		h.state[Slot] = hdr.slot
	}
}

func (h *Heap) buddyMerge() {
	log.Printf("buddyMerge() %v\n", h)
	loc := h.state[Loc]
	slot := h.state[Slot]
	bdy := h.state[Bdy]
	h.resetState()

	front := uintMin(loc, bdy)
	// remove them both, insert them merge
	h.removeCell(loc)
	h.removeCell(bdy)
	h.insertCell(front, slot+1)

	// Check again if you can merge
	h.state[Type] = BuddyChk
	h.state[Loc] = front
	h.state[Bdy] = h.getBuddy(front)
}

// Just an indirection to show that buddy merge failed
func (h *Heap) buddyFail() {
	log.Printf("buddyFail() %v\n", h)
	h.resetState()
	h.state[Type] = Idle
}
