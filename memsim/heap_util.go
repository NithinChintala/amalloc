package memsim

import (
	"log"
)

func (h *Heap) removeCell(loc uint) {
	log.Printf("removeCell(%d) %v\n", loc, h)
	oldCell := h.readCell(loc)
	idx := slotToIdx(oldCell.slot)

	if oldCell.prev == NullPtr {
	// Front of free list
		if oldCell.next == NullPtr {
		// singelton list
			h.heads[idx] = NullPtr
		} else {
		// front of list len > 1
			newFront := h.readCell(oldCell.next)
			newFront.prev = NullPtr
			h.writeCell(oldCell.next, newFront)
			h.heads[idx] = oldCell.next
		}
	} else if oldCell.next == NullPtr {
	// At end of list len > 1
		prevCell := h.readCell(oldCell.prev)
		prevCell.next = oldCell.next
		h.writeCell(oldCell.prev, prevCell)
	} else {
	/// At Middle of list
		prevCell := h.readCell(oldCell.prev)
		nextCell := h.readCell(oldCell.next)

		prevCell.next = oldCell.next
		nextCell.prev = oldCell.prev

		h.writeCell(oldCell.prev, prevCell)
		h.writeCell(nextCell.prev, nextCell)
	}
}

func (h *Heap) insertCell(loc, slot uint) {
	log.Printf("insertCell(loc=%d, slot=%d) %v\n", loc, slot, h)
	idx := slotToIdx(slot)
	newCell := Cell{}
	newCell.slot = slot
	newCell.used = false

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
}

// Transfer h.state to h.prevState
// Deletes everything from h.state
func (h *Heap) resetState() {
	for k := range h.prevState {
		delete(h.prevState, k)
	}
	for k, v := range h.state {
		h.prevState[k] = v
		delete(h.state, k)
	}
}

func (h *Heap) getBuddy(loc uint) uint {
	hdr := h.readHeader(loc)
	if hdr.slot >= MaxPwr {
	// MaxPwr slot has not buddies
		return NullPtr
	}
	return loc ^ (1 << hdr.slot)
}

func (h *Heap) getPrevState() string {
	switch h.prevState[Type] {
	case Idle:
		return "Idle"
	case Split:
		return "Split"
	case CheckAvail:
		return "Check Avail"
	case SetHead:
		return "Set Head"
	case ValSet:
		return "Val Set"
	case FreeSet:
		return "Free Set"
	case BuddyChk:
		return "Buddy Check"
	case BuddyFail:
		return "Buddy Fail"
	case BuddyMerge:
		return "Buddy Merge"
	case OutOfMem:
		return "Out of Memory"
	default:
		return ""
	}
}
