package memsim

import (
	"fmt"
	"strings"
	"strconv"
	_ "time"
)

const (
	Clear = "\033[2J\033[0;0H"
)

var (
	Blank = strings.Repeat(" ", 8)
)

func Render(h *Heap) {
	// The first Byte MUST be a header / cell
	mat := make([][]string, MaxPwr)
	var heapSize uint = 1 << MaxPwr
	for i := range mat {
		mat[i] = make([]string, heapSize)
	}

	var curr uint = 0
	var idx uint
	var i uint
	var hdr Header
	for curr < heapSize {
		hdr = h.readHeader(curr)
		idx = slotToIdx(hdr.slot)
		for i = 0; i < 1 << hdr.slot; i++ {
			mat[idx][curr + i] = byteString(h.mem[curr + i])
		}
		curr += 1 << hdr.slot
	}
	fmt.Print(Clear)
	row := make([]string, heapSize)
	for i := MaxPwr - 1; i > -1; i-- {
		for j := range mat[i] {
			if mat[i][j] == "" {
				mat[i][j] = Blank
			}
			row[j] = mat[i][j]
		}
		fmt.Printf("%2d: %s\n", 1 << idxToSlot(uint(i)), strings.Join(row, " "))
	}
	/*
	builder := make([]string, 1 << MaxPwr)
	for i, b := range h.mem {
		builder[i] = fmt.Sprintf("%08s", strconv.FormatInt(int64(b), 2))
	}

	fmt.Printf("%s%s\n", Clear, strings.Join(builder, " "))
	*/
}

func byteString(b byte) string {
	return fmt.Sprintf("%08s", strconv.FormatInt(int64(b), 2))
}

func Anim(h *Heap) {
	var cmd string
	count := 0
	for {
		Render(h)
		fmt.Printf("\n%d %v\n>>> ", count, h)
		h.Step()
		count++
		fmt.Scanln(&cmd)
		//time.Sleep(1 * time.Second)
	}
}