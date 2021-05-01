package memsim

import (
	"fmt"
	"strings"
	"time"
	"regexp"
	"bufio"
	"os"

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
	PointUp = "↑" + strings.Repeat(" ", 7)
	Check = color.Magenta("???" + strings.Repeat(" ", 5))
	Fail = color.Red("xxx" + strings.Repeat(" ", 5))
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

	for char, loc := range h.vars {
		hdr := h.readHeader(loc - 1)
		mat[slotToAnimRow(hdr.slot) + 1][loc + LPad] = PointUp
		mat[slotToAnimRow(hdr.slot) + 2][loc + LPad] = fmt.Sprintf("%c", char) + strings.Repeat(" ", 7)
	}
}

// Set anything specific related to h.state
func setState(h *Heap, mat [][]string) {
	switch h.prevState[Type] {
	case Split: 
		slot := h.prevState[Slot]
		row := slotToAnimRow(slot)
		col := h.prevState[Loc] + LPad
		mat[row][col] = PointDown
		mat[row][col + (1 << (slot - 1))] = PointDown
	case SetHead:
		fallthrough
	case FreeSet:
		slot := h.prevState[Slot]
		row := slotToAnimRow(slot) - 1
		col := h.prevState[Loc] + LPad
		mat[row][col] = SetByte
	case BuddyChk:
		var row, col uint
		loc := h.prevState[Loc]
		hdr := h.readHeader(loc)
		row = slotToAnimRow(hdr.slot) - 1
		col = loc + LPad
		mat[row][col] = Check

		bdy := h.prevState[Bdy]
		if bdy != NullPtr {
			buddy := h.readHeader(bdy)
			row = slotToAnimRow(buddy.slot) - 1
			col = bdy + LPad
			mat[row][col] = Check
		}
	case BuddyMerge:
		loc := h.prevState[Loc]
		slot := h.prevState[Slot]
		bdy := h.prevState[Bdy]

		row := slotToAnimRow(slot)
		mat[row][loc + LPad] = color.Green(PointUp)
		mat[row][bdy + LPad] = color.Green(PointUp)

	case BuddyFail:
		var row, col uint
		loc := h.prevState[Loc]
		hdr := h.readHeader(loc)
		row = slotToAnimRow(hdr.slot) - 1
		col = loc + LPad
		mat[row][col] = Fail

		bdy := h.prevState[Bdy]
		buddy := h.readHeader(bdy)
		row = slotToAnimRow(buddy.slot) - 1
		col = bdy + LPad
		mat[row][col] = Fail
	}
}

func Anim(h *Heap) {
	count := 0
	fmt.Print(ClearOrigin)
	mallocRegex := regexp.MustCompile(`^([A-Za-z]{1}) = malloc\(([1-9]{1}[0-9]*)\)`)
	freeRegex := regexp.MustCompile(`^free\(([A-Za-z]{1})\)`)
	setValRegex := regexp.MustCompile(`^([A-Za-z]{1}) = ([1-9]{1}[0-9]*)`)
	getValRegex := regexp.MustCompile(`^([A-Za-z]{1})`)
	reader := bufio.NewReader(os.Stdin)
	for {
		Render(h)
		fmt.Printf("\n%s%d %v\n>>> ", ClearEnd, count, h)
		count++
		if h.prevState[Type] == Idle && h.state[Type] == Idle {
		// Wait for a user command
			cmd, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Issue when reading input")
				os.Exit(1)
			}
			cmd = strings.TrimSuffix(cmd, "\r\n")
			if mallocRegex.MatchString(cmd) {
			// <var> = malloc(<size>)
				parsed := mallocRegex.FindStringSubmatch(cmd)
				h.Malloc(parsed[1], mustAtoui(parsed[2]))
			} else if freeRegex.MatchString(cmd) {
			// free(<var>)
				parsed := freeRegex.FindStringSubmatch(cmd)
				h.Free(parsed[1])
				h.Step()
			} else if setValRegex.MatchString(cmd) {
			// <var> = <val>					
				fmt.Println(cmd)
				os.Exit(0)
			} else if getValRegex.MatchString(cmd) {
			// <var>
				fmt.Println(cmd)
				os.Exit(0)
			} else {
			// Bad syntax
				fmt.Printf("Bad command syntax: '%s'\n", cmd)
				os.Exit(1)
			}
		} else {
		// Wait for a second and continue
			//time.Sleep(600 * time.Millisecond)
			time.Sleep(1 * time.Second)
			//reader.ReadString('\n')
			h.Step()
		}
	}
}
