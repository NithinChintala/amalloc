package main

import (
	"github.com/NithinChintala/amalloc/memsim"
	"log"
	"fmt"
)

type DevNull struct {
}

func (dn *DevNull) Write(p []byte) (n int, err error) {
	return 0, nil
}
func main() {
	interact()
	//debug()
}

func debug() {
	log.SetFlags(log.Flags() & ^(log.Ldate | log.Ltime))
	h := memsim.NewHeap()
	h.Malloc("x", 1)
	for i := 0; i < 8; i++ {
		fmt.Printf("%d ", i)
		h.Step()
	}
	h.Free(1)
	for i := 0; i < 8; i++ {
		fmt.Printf("%d ", i)
		h.Step()
	}
}

func interact() {
	log.SetOutput(&DevNull{})
	h := memsim.NewHeap()
	memsim.Anim(h)
}