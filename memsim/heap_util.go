package memsim

import (
	"log"
)

// Update the pointers in the cell free list
func (h *Heap) removeCell(loc uint) {
	log.Printf("removeCell(%d) %v\n", loc, h)
	oldCell := h.readCell(loc)
	idx := slotToIdx(oldCell.slot)

	if oldCell.prev == NullPtr {
	// Front of free list
		if oldCell.next == NullPtr {
		// singelton list
			//log.Println("singleton list")
			h.heads[idx] = NullPtr
		} else {
		// front of list len > 1
			//log.Println("front, len > 1")
			//fmt.Println(oldCell)
			newFront := h.readCell(oldCell.next)
			newFront.prev = NullPtr
			h.writeCell(oldCell.next, newFront)
			h.heads[idx] = oldCell.next
		}
	} else if oldCell.next == NullPtr {
	// At end of list len > 1
		//log.Println("end, len > 1")
		prevCell := h.readCell(oldCell.prev)
		prevCell.next = oldCell.next
		h.writeCell(oldCell.prev, prevCell)
	} else {
	/// At Middle of list
		//log.Println("middle")
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
	//log.Printf("after removeCell(%d) %v\n", loc, h)
}

func (h *Heap) insertCell(loc, slot uint, used bool) {
	log.Printf("insertCell(loc=%d, slot=%d, used=%t) %v\n", loc, slot, used, h)
	idx := slotToIdx(slot)
	newCell := Cell{}
	newCell.slot = slot
	newCell.used = used
	
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
	//log.Printf("after insertCell(loc=%d, slot=%d, used=%t) %v\n", loc, slot, used, h)
}

// Transfered h.state to h.prevState
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
