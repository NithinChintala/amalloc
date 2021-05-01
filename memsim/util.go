package memsim

import (
	"fmt"
	"strconv"
)

// The minimum power of 2 that is >= x
func minPower2(x uint) uint {
	var power uint = 0
	var curr uint = 1
	for curr < x {
		power++
		curr <<= 1
	}
	return power
}

// Converts the slot to the index in the slice
func slotToIdx(slot uint) uint {
	return slot - MinPwr
}

// Converts the index of something in the slice to
// the appropriate slot
func idxToSlot(idx uint) uint {
	return idx + MinPwr
}

// Returns the minimum of the two uints
func uintMin(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

// Converts a byte to a boolean
// b > 0 -> true
// b = 0 -> flase
func byteToBool(b byte) bool {
	return b > 0
}

// Converts a boolean to a byte
// true  -> 0b00000001
// false -> 0b00000000
func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func strToUint(s string) uint {
	return uint([]rune(s)[0])
}

// Have the slots decrease from 16 -> 2
func slotToAnimRow(slot uint) uint {
	return TPad + MaxPwr - slotToIdx(slot) - 1
}

func animRowToSlot(row uint) uint {
	return idxToSlot(TPad + MaxPwr - row - 1)
}

func byteString(b byte) string {
	return fmt.Sprintf("%08s", strconv.FormatInt(int64(b), 2))
}

func numPad8(n uint) string {
	return fmt.Sprintf("%-8d", n)
}

// Converts the strings to an unsigned int, panics otherwise
func mustAtoui(s string) uint {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return uint(n)
}
