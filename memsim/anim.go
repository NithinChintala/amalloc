package memsim

import (
	"fmt"
	"strings"
	"strconv"
	"time"

	"github.com/NithinChintala/amalloc/color"
)

const (
	Origin = "\033[0;0H"
	ClearOrigin = "\033[2J\033[0;0H"
	ClearLine = "\033[2K"
	ClearEnd = "\033[0J"

	// Left, Top, Bottom padding
	LPadLen = 4
	LPad = 1
	TPad = 2
	BPad = 2
)

var (
	BlankByte = strings.Repeat(" ", 8)
	BlankAnno = strings.Repeat(" ", 3)
	SetByte = strings.Repeat("v", 8)
	PointDown = "↓" + strings.Repeat(" ", 7)
)

// Format:
// 	- The output is 17 col x 8 row
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

	// Set everything to blank initially, fill stuff in as you go
	for row := range mat {
		mat[row][0] = BlankAnno
		for col := LPad; col < len(mat[row]); col++ {
			mat[row][col] = BlankByte
		}
	}

	setMem(h, mat)
	setAnnotate(h, mat)
	setState(h, mat)

	// Join cols by spaces, rows by new line
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
}

// Set all the annotations in the render. These include:
// 	- slot sizes on the left, i.e. " 8:"
// 	- pointers to currently allocated variables
//	- pointers to the heads (if not NilPtr) of each slot
func setAnnotate(h *Heap, mat [][]string) {
	// Add the slot sizes on the left
	var i uint = TPad
	var checking, willSplit, found bool
	var slot uint
	var str string
	for ; i < TPad + MaxPwr; i++ {
		slot = animRowToSlot(i)
		str = fmt.Sprintf("%2d:", 1 << slot)

		checking = h.state[Type] == CheckAvail && h.state[Slot] == slot
		willSplit = h.state[Type] == Split && h.prevState[Type] == CheckAvail && h.prevState[Slot] == slot
		found = h.state[Type] == SetHead && h.prevState[Type] == CheckAvail && h.prevState[Slot] == slot
		if checking {
		// Magenta for checking if this slot has something
			str = color.Magenta(str)
		} else if willSplit || found {
		// Green for when after checking, there is something
			str = color.Green(str)
		}

		mat[i][0] = str
	}

	for i := 0; i < MaxPwr; i++ {
		slot := idxToSlot(uint(i))
		if loc := h.heads[i]; loc != NullPtr {
			mat[slotToAnimRow(slot) - 1][loc + LPad] = PointDown
			mat[slotToAnimRow(slot) - 2][loc + LPad] = numPad8(1 << slot)
		}
	}
}

// Set anything specific related to h.state
func setState(h *Heap, mat [][]string) {
	if h.prevState[Type] == Split {
		slot := h.prevState[Slot]
		row := slotToAnimRow(slot)
		col := h.prevState[Loc] + LPad
		mat[row][col] = PointDown
		mat[row][col + (1 << (slot - 1))] = PointDown
	} else if h.prevState[Type] == SetHead {
		slot := h.prevState[Slot]
		row := slotToAnimRow(slot) - 1
		col := h.prevState[Loc] + LPad
		mat[row][col] = SetByte
	}
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

func Anim3(h *Heap) {
	var cmd string
	count := 0
	fmt.Print(ClearOrigin)
	for {
		Render(h)
		fmt.Printf("\n%s%d %v\n>>> ", ClearEnd, count, h)
		h.Step()
		count++
		fmt.Scanln(&cmd)
	}
}

func Anim(h *Heap) {
	var cmd string
	count := 0
	fmt.Print(ClearOrigin)
	for {
		Render(h)
		fmt.Printf("\n%s%d %v\n>>> ", ClearEnd, count, h)
		h.Step()
		count++
		//fmt.Scanln(&cmd)
		if h.prevState[Type] == Idle && h.state[Type] == Idle {
			// Step one more time to Idle
			time.Sleep(1 * time.Second)
			Render(h)
			fmt.Printf("\n%s%d %v\n>>> ", ClearEnd, count, h)
			fmt.Scanln(&cmd)
			n, err := strconv.Atoi(cmd)
			if err != nil {
				panic(err)
			}
			h.Malloc(uint(n))
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func Anim2(h *Heap) {
	var cmd string
	count := 0
	fmt.Print(ClearOrigin)
	for {
		Render(h)
		fmt.Printf("\n%s%d %v\n>>> ", ClearEnd, count, h)
		h.Step()
		count++
		if h.prevState[Type] == Idle {
			Render(h)
			fmt.Scanln(&cmd)
			n, err := strconv.Atoi(cmd)
			if err != nil {
				panic(err)
			}
			h.Malloc(uint(n))
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}