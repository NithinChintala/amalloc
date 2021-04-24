package main

import (
	"fmt"
	"github.com/NithinChintala/ascii-malloc/color"
)

func main() {
	fmt.Println(color.Red("hello ") + color.Magenta("world!"))
}