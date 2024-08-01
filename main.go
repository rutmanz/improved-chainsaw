package main

import (
	"fmt"

	"github.com/rutmanz/improved-chainsaw/parser"
)

func main() {
	ctx := parser.Parse("3(a|b)+\\.{2,5}d")
	for _, token := range ctx.Tokens {
		fmt.Println(token.ToString(""))
	}

}
