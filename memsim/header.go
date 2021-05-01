package memsim

type Cell struct {
	used bool
	slot uint
	next uint
	prev uint
}

type Header struct {
	used bool
	slot uint
}

// Reads the cell at the given memory location
func (h *Heap) readCell(loc uint) Cell {
	byte1 := h.mem[loc]
	byte2 := h.mem[loc+1]

	c := Cell{}
	c.used = byteToBool(byte1 & (1 << 7))
	c.slot = idxToSlot(uint((byte1 & 0b01100000) >> 5))
	c.prev = uint(byte1 & 0b00011111)
	c.next = uint(byte2)

	return c
}

// Write the given cell at the given memory location
func (h *Heap) writeCell(loc uint, c Cell) {
	var byte1 byte = boolToByte(c.used) << 7
	byte1 |= byte(slotToIdx(c.slot) << 5)
	byte1 |= byte(c.prev)

	byte2 := byte(c.next)
	h.mem[loc] = byte1
	h.mem[loc+1] = byte2
}

// Reads the header at the given memory location
func (h *Heap) readHeader(loc uint) Header {
	byte1 := h.mem[loc]

	hdr := Header{}
	hdr.used = byteToBool(byte1 & (1 << 7))
	hdr.slot = idxToSlot(uint((byte1 & 0b01100000) >> 5))

	return hdr
}

// Write the given header at the given memory location
func (h *Heap) writeHeader(loc uint, hdr Header) {
	var byte1 byte = boolToByte(hdr.used) << 7
	byte1 |= byte(slotToIdx(hdr.slot) << 5)

	h.mem[loc] = byte1
}
