package main

import (
	"github.com/NithinChintala/ascii-malloc/memsim"
	"log"
	"fmt"
)

func main() {
	log.SetFlags(log.Flags() & ^log.Ldate)
	h := memsim.NewHeap(4)
	h.Malloc(4)
	h.Malloc(8)
	fmt.Println(h)
	h.Malloc(64)
}