package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/laupse/dagger-demo/math"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		panic("No enough args")
	}
	a, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(args[2])
	if err != nil {
		panic(err)
	}
	result := math.Add(a)
	if _, err := fmt.Printf("%d + %d = %d", a, b, result); err != nil {
		panic(err)
	}
}
