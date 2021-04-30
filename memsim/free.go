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
		h.state[Bdy] = bdy
	} else {
	// Buddy is not free
		h.state[Type] = BuddyFail
		h.state[Loc] = loc
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

	// remove the buddy, insert them both
	h.removeCell(bdy)
	h.insertCell(front, slot + 1, false)

	// Check again if you can merge
	h.state[Type] = BuddyChk
	h.state[Loc] = loc
	h.state[Bdy] = h.getBuddy(front)
}

func (h *Heap) buddyFail() {
	// No buddy available, so finally insert the cell in
	log.Printf("buddyFail() %v\n", h)
	loc := h.state[Loc]
	slot := h.state[Slot]
	h.resetState()

	h.insertCell(loc, slot, false)

	h.state[Type] = Idle
}