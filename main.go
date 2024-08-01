package main

import (
	"fmt"

	"github.com/rutmanz/improved-chainsaw/parser"
)

func main() {
	matched := parser.Match("a|b", "a")
	fmt.Printf("Matched: %t\n", matched)
}
