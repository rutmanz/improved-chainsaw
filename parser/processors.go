package parser

import (
	"strconv"
)

func (c *Context) process() {
	input := c.getRegex()
	// fmt.Println(input)
	// fmt.Println(strings.Repeat(" ", c.index) + "^")
	switch input[c.index] {
	case '(':
		c.processGroup()
	case '[':
		c.processCharset()
	case '|':
		c.processAlternation()
	case '?':
		c.processShorthandQuantifier(0, 1)
	case '*':
		c.processShorthandQuantifier(0, -1)
	case '+':
		c.processShorthandQuantifier(1, -1)
	case '{':
		c.processExplicitQuantifier()
	case '.':
		c.push(TokenCharsetAny{})
	case '\\':
		c.processEscape()
	default:
		c.push(TokenLiteral{rune(input[c.index])})
	
	}
}

var (
	// Charsets
	charsetDigits interfaceCharset = TokenCharsetRange{nil, '0', '9'}
	charsetNotDigits interfaceCharset = TokenCharsetNot{nil, charsetDigits}
	charsetWhitespace interfaceCharset = TokenCharsetLiterals{nil, map[byte]struct{}{'\t': {}, '\n': {}, '\v': {}, '\f': {}, '\r': {}, ' ': {}}}
	charsetNotWhitespace interfaceCharset = TokenCharsetNot{nil, charsetWhitespace}
	charsetWord interfaceCharset = TokenCharsetLiterals{nil, map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}, 'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {}, 'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {}, 'm': {}, 'n': {}, 'o': {}, 'p': {}, 'q': {}, 'r': {}, 's': {}, 't': {}, 'u': {}, 'v': {}, 'w': {}, 'x': {}, 'y': {}, 'z': {}, 'A': {}, 'B': {}, 'C': {}, 'D': {}, 'E': {}, 'F': {}, 'G': {}, 'H': {}, 'I': {}, 'J': {}, 'K': {}, 'L': {}, 'M': {}, 'N': {}, 'O': {}, 'P': {}, 'Q': {}, 'R': {}, 'S': {}, 'T': {}, 'U': {}, 'V': {}, 'W': {}, 'X': {}, 'Y': {}, 'Z': {}}}
	charsetNotWord interfaceCharset = TokenCharsetNot{nil, charsetWord}
)

func (c *Context) processEscape() {
	input := c.getRegex()
	c.index++
	switch input[c.index] {
	case 'd':
		c.push(charsetDigits)
	case 'D':
		c.push(charsetNotDigits)
	case 's':
		c.push(charsetWhitespace)
	case 'S':
		c.push(charsetNotWhitespace)
	case 'w':
		c.push(charsetWord)
	case 'W':
		c.push(charsetNotWord)
	default:
		c.push(TokenLiteral{rune(input[c.index])})
	}
}

func (c *Context) processGroup() {
	groupCtx := &Context{index: c.index + 1, regex: c.regex, Tokens: make([]token, 0, len(c.regex))}
	for c.regex[groupCtx.index] != ')' {
		groupCtx.process()
		groupCtx.index++
	}
	c.index = groupCtx.index
	c.push(TokenGroupCaptured{groupCtx.Tokens})
}

type charitem struct {
	start byte
	end   byte
}

func (c *Context) processCharset() {
	input := c.getRegex()
	items := make([]charitem, 0)
	c.index++
	for input[c.index] != ']' {
		switch input[c.index] {
		case '-':
			items[len(items)-1].end = input[c.index+1]
			c.index++
		default:
			items = append(items, charitem{start: input[c.index]})
			c.index++
		}
	}
	charSet := make(map[byte]struct{})
	for _, item := range items {
		if item.end == 0x00 {
			charSet[item.start] = struct{}{}
		} else {
			for i := item.start; i <= item.end; i++ {
				charSet[i] = struct{}{}
			}
		}
	}
	c.push(TokenCharsetLiterals{nil, charSet})
}

func (c *Context) processAlternation() {
	input := c.getRegex()
	rightSideContext := &Context{index: c.index + 1, regex: c.regex, Tokens: make([]token, 0)}
	for rightSideContext.index < len(input) && input[rightSideContext.index] != ')' {
		rightSideContext.process()
		rightSideContext.index++
	}

	leftToken := TokenGroupUncaptured{c.Tokens} // All existing tokens in current context are on the left side of the alternation
	rightToken := TokenGroupUncaptured{rightSideContext.Tokens}

	c.index = rightSideContext.index - 1       // -1 because the group will increment the index
	c.Tokens = make([]token, 0, cap(c.Tokens)) // Clear old tokens
	c.push(TokenAlternate{leftToken, rightToken})
}

func (c *Context) processShorthandQuantifier(min int, max int) {
	prevToken := c.Tokens[len(c.Tokens)-1]
	c.Tokens[len(c.Tokens)-1] = TokenQuantifier{ // Replace the last token with a quantifier
		min:   min,
		max:   max,
		token: prevToken,
	}
}

func (c *Context) processExplicitQuantifier() {
	input := c.getRegex()
	var min, max int
	var minChars, maxChars []byte = make([]byte, 0), make([]byte, 0)
	hasFoundComma := false
	c.index++
	for input[c.index] != '}' {
		if input[c.index] == ',' {
			hasFoundComma = true
			c.index++
			continue
		}
		if hasFoundComma {
			maxChars = append(maxChars, input[c.index])
		} else {
			minChars = append(minChars, input[c.index])
		}
		c.index++
	}

	if !hasFoundComma {
		min, _ = strconv.Atoi(string(minChars))
		max = min
	} else {
		if len(minChars) == 0 {
			min = 0
		} else {
			min, _ = strconv.Atoi(string(minChars))
		}
		if len(maxChars) == 0 {
			max = -1
		} else {
			max, _ = strconv.Atoi(string(maxChars))
		}
	}
	c.processShorthandQuantifier(min, max)
}
