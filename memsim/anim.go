package memsim

import (
	"fmt"
	"strings"
)

func Render(h *Heap) string {
	mat := make([][]string, MaxPwr)
	for i := range mat {
		mat[i] = make([]string, 1 << MaxPwr)
	}

	for _, b := range h.mem {
		fmt.Printf("%d ", b)
	}

	out := make([]string, MaxPwr)
	for _, row := range mat {
		out = append(out, strings.Join(row, " "))
	}

	return strings.Join(out, "\n")
}