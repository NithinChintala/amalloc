package main

import (
	"github.com/NithinChintala/ascii-malloc/memsim"
	"log"
)

func main() {
	log.SetFlags(log.Flags() & ^log.Ldate)
	h := memsim.NewHeap()
	h.Malloc(1)
	h.Step()
	h.Step()
	h.Step()
}