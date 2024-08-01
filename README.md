# Improved Chainsaw
- Not a chainsaw
- Not an improvement

## Features
This is a simple golang regex parser and matcher using Thompson’s construction. It only supports matching the full input, and will not match substrings. It supports:
- Character class escapes (`\w`, `\d`, `\s`) and their negations (`\W`, `\D`, `\S`) 
- Character classes (`[A-Z]`, `[abc]`)
- Wildcard symbols (`.`)
- Quantifers (`*`, `+`, `?`, `{n}`, `{n,}`, `{n,m}`)
- Groups (`(abc)`)
- Alternation (`a|b`, `a|b|c`)
- Character Literals (`a`, `b`, `c`)

## Usage

```go
package main

import (
	"fmt"

	"github.com/rutmanz/improved-chainsaw/parser"
)

func main() {
	matched := parser.Match("a|b", "a")
	fmt.Printf("Matched: %t\n", matched)
}

```