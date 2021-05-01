package main

import (
	"fmt"
	"log"
	"flag"

	"github.com/NithinChintala/amalloc/memsim"
)

type DevNull struct{}
func (dn *DevNull) Write(p []byte) (n int, err error) { return 0, nil }

func main() {
	
	speed := getSpeed()
	//interact()
	//debug()
}

func getSpeed() string {
	validSpeeds := []string{"step", "slow", "norm", "fast", "inst"}
	speedPtr := flag.String("speed", "norm", "The speed of the animation")
	flag.Parse()

	for _, speed := range validSpeeds {
		if *speedPtr == speed {
			return speed
		}
	}
	log.Fatal("Invalid command line arguments")
}

func debug() {
	log.SetFlags(log.Flags() & ^(log.Ldate | log.Ltime))
	h := memsim.NewHeap()
	h.Malloc("x", 1)
	for i := 0; i < 8; i++ {
		fmt.Printf("%d ", i)
		h.Step()
	}
	h.Free("x")
	for i := 0; i < 15; i++ {
		fmt.Printf("%d ", i)
		h.Step()
	}
}

func interact() {
	log.SetOutput(&DevNull{})
	h := memsim.NewHeap()
	memsim.Anim(h)
}
