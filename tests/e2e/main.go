package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	numStr := os.Getenv("NUM")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		os.Exit(1)
	}
	if num%2 != 0 {
		fmt.Printf("odd. num=%d\n", num)
		os.Exit(1)
	}
	fmt.Printf("even. num=%d\n", num)
}
