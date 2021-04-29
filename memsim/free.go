package memsim

import (
	"log"
)

func (h *Heap) freeSet() {
	log.Printf("freeSet() %v\n", h)
	loc := h.state[Loc]
	h.resetState()

	hdr := h.readHeader(loc)
	hdr.used = false
	h.writeHeader(loc, hdr)

	bdyLoc := h.getBuddy(loc)
	if bdyLoc == NullPtr {
	// Only when slot is 4
		h.insertCell(loc, hdr.slot, false)
		h.state[Type] = Idle
		return
	}

	h.state[Type] = BuddyChk
	h.state[Loc] = loc
	h.state[Bdy] = bdyLoc
}

func (h *Heap) buddyChk() {
	log.Printf("buddyChk() %v\n", h)
	loc := h.state[Loc]
	// bdy should not be NullPtr
	bdy := h.state[Bdy]
	h.resetState()

	hdr := h.readHeader(loc)
	buddy := h.readHeader(bdy)
	if !buddy.used && buddy.slot == hdr.slot {
	// Found the buddy, has same size + is not used
		h.state[Type] = BuddyMerge
		h.state[Loc] = loc
		h.state[Slot] = hdr.slot
	} else {
	// Buddy is not free
		h.state[Type] = BuddyMerge
		h.state[Loc] = loc
		h.state[Slot] = hdr.slot
	}
}
