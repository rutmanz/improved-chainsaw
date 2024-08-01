package main

import (
	"fmt"

	"github.com/rutmanz/improved-chainsaw/parser"
)

func main() {
	fmt.Printf("\nMATCH: %t\n", parser.Match("(a|b)c", "bc"))
}
