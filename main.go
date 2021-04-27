package main

import (
	"github.com/NithinChintala/ascii-malloc/memsim"
	"log"
)

type DevNull struct {
}

func (dn *DevNull) Write(p []byte) (n int, err error) {
	return 0, nil
}
func main() {
	interact()
}

func debug() {
	log.SetFlags(log.Flags() & ^(log.Ldate | log.Ltime))
	h := memsim.NewHeap()
	h.Malloc(1)
	for {
		h.Step()
	}
}

func interact() {
	log.SetOutput(&DevNull{})
	h := memsim.NewHeap()
	h.Malloc(1)
	memsim.Anim(h)
}