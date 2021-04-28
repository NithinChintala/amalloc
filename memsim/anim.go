package memsim

import (
	"fmt"
	"strings"
	"strconv"
	_ "time"
)

const (
	Clear = "\033[2J\033[0;0H"
	Origin = "\033[0;0H"
	//ClearLine = "\033[1K"
	ClearLine = "\033[2K"
	ClearEnd = "\033[0J"
	// Left, Top, Bottom padding
	LPadLen = 4
	LPad = 1
	TPad = 2
	BPad = 2
)

var (
	Blank = strings.Repeat(" ", 8)
	SetByte = strings.Repeat("v", 8)
)

// Format:
// 	- The output is a 17 col x 8 row
// 	- The first col is 3 characters long 
// 	- the next 16 cols are 8 characters long 
// 	- All columns are space seperated
// 	- All rows are new line seperated
// 	- Below is an example format
//   __________________
// 0|                  |
// 1|                  |
// 2|16: xx xx ..... xn|
// 3| 8: xx xx ..... xn|
// 4| 4: xx xx ..... xn|
// 5| 2: xx xx ..... xn|
// 6|                  |
// 7|__________________|
//     0  1  2 ..... 16

// Renders a the current state of the heap
func Render(h *Heap) {
	// Initialize the rendering matrix
	mat := make([][]string, TPad + MaxPwr + BPad)
	for i := range mat {
		mat[i] = make([]string, LPad + (1 << MaxPwr))
	}

	setMem(h, mat)
	setAnnotate(h, mat)
	setState(h, mat)

	rows := make([]string, TPad + MaxPwr + BPad)
	for i := range rows {
		rows[i] = strings.Join(mat[i], " ")
	}
	out := strings.Join(rows, "\n")
	fmt.Printf("%s%s\n", Origin, out)
}

// Set all the bytes of memory in the correct row and col
func setMem(h *Heap, mat [][]string) {
	// The first byte in h.mem MUST be a header / cell
	// It doesn't matter which all we need is the slot
	var curr, heapSize uint
	var hdr Header
	heapSize = 1 << MaxPwr
	for curr < heapSize {
		hdr = h.readHeader(curr)
		var i, row, col uint
		for ; i < 1 << hdr.slot; i++ {
			row = slotToAnimRow(hdr.slot)
			col = curr + i + LPad
			mat[row][col] = byteString(h.mem[curr + i])
		}
		curr += 1 << hdr.slot
	}
	
	// Set everything unsets in memory to 8 spaces
	for i := MaxPwr - 1; i > -1; i-- {
		var j, row, col uint
		for ; j < heapSize; j++ {
			row = TPad + uint(i)
			col = LPad + j
			if mat[row][col] == "" {
				mat[row][col] = Blank
			}
		}
	}
}

// Set all the annotations in the render. These include:
// 	- slot sizes on the left, i.e. " 8:"
// 	- pointers to currently allocated variables
//	- pointers to the heads (if not NilPtr) of each slot
func setAnnotate(h *Heap, mat [][]string) {
	// Add the slot sizes on the left
	var i uint = TPad
	for ; i < TPad + MaxPwr; i++ {
		slot := animRowToSlot(i)
		mat[i][0] = fmt.Sprintf("%2d:", 1 << slot)
	}
}

// Set anything specific related to h.state
func setState(h *Heap, mat [][]string) {

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

func Anim(h *Heap) {
	var cmd string
	count := 0
	fmt.Print(Clear)
	for {
		Render(h)
		fmt.Printf("\n%s%d %v\n>>> ", ClearEnd, count, h)
		h.Step()
		count++
		fmt.Scanln(&cmd)
		//time.Sleep(1 * time.Second)
	}
}