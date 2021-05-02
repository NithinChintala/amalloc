package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/NithinChintala/amalloc/memsim"
)

type DevNull struct{}

func (dn *DevNull) Write(p []byte) (n int, err error) { return 0, nil }

func main() {
	interact()
}

func interact() {
	speed, err := getSpeed()
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(&DevNull{})
	h := memsim.NewHeap()
	waitFunc := getSpeedFunc(speed)
	memsim.Anim(h, waitFunc)
}

// Hard coded debugging
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

func getSpeed() (string, error) {
	validSpeeds := []string{"step", "slow", "norm", "fast", "inst"}
	speedPtr := flag.String("speed", "norm", "The speed of the animation: step|slow|norm|fast|inst")
	flag.Parse()

	for _, speed := range validSpeeds {
		if *speedPtr == speed {
			return speed, nil
		}
	}
	return "", fmt.Errorf("Invalid speed argument: '%s'", *speedPtr)
}

func getSpeedFunc(speed string) func() {
	switch speed {
	case "step":
		return func() { fmt.Scanln() }
	case "slow":
		return func() { time.Sleep(1500 * time.Millisecond) }
	case "norm":
		return func() { time.Sleep(1 * time.Second) }
	case "fast":
		return func() { time.Sleep(500 * time.Millisecond) }
	case "inst":
		return func() {}
	default:
		panic(fmt.Sprintf("getSpeedFunc(%s) is invalid\n", speed))
	}
}
