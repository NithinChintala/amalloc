package memsim

func minPower2(x uint) uint {
	var power uint = 0
	var curr uint = 1
	for curr < x {
		power++
		curr <<= 1
	}
	return power
}

func slotToIdx(slot uint) uint {
	return slot - MinPwr
}

func idxToSlot(idx uint) uint {
	return idx + MinPwr
}

func uintMin(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

func byteToBool(b byte) bool {
	return b > 0
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}