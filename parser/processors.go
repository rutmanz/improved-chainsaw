package parser

import (
	"strconv"
)

func (c *Context) parseNext() {
	input := c.getRegex()
	// fmt.Println(input)
	// fmt.Println(strings.Repeat(" ", c.index) + "^")
	switch input[c.index] {
	case '(':
		c.parseGroup()
	case '[':
		c.parseCharset()
	case '|':
		c.parseAlternation()
	case '?':
		c.parseShorthandQuantifier(0, 1)
	case '*':
		c.parseShorthandQuantifier(0, -1)
	case '+':
		c.parseShorthandQuantifier(1, -1)
	case '{':
		c.parseExplicitQuantifier()
	case '.':
		c.push(charsetAny)
	case '\\':
		c.parseEscape()
	default:
		c.push(TokenLiteral{input[c.index]})

	}
}

var (
	// Charsets
	charsetDigits        interfaceCharset = TokenCharsetRange{nil, '0', '9'}
	charsetNotDigits     interfaceCharset = TokenCharsetNot{nil, charsetDigits}
	charsetWhitespace    interfaceCharset = TokenCharsetLiterals{nil, map[byte]struct{}{'\t': {}, '\n': {}, '\v': {}, '\f': {}, '\r': {}, ' ': {}}}
	charsetNotWhitespace interfaceCharset = TokenCharsetNot{nil, charsetWhitespace}
	charsetWord          interfaceCharset = TokenCharsetCompound{nil, []interfaceCharset{charsetDigits, TokenCharsetRange{nil, 'a', 'z'}, TokenCharsetRange{nil, 'A', 'Z'}, TokenCharsetLiterals{nil, map[byte]struct{}{'_': {}}}}}
	charsetNotWord       interfaceCharset = TokenCharsetNot{nil, charsetWord}
	charsetAny           interfaceCharset = TokenCharsetRange{nil, 0x00, 0xFF}
)

func (c *Context) parseEscape() {
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
		c.push(TokenLiteral{input[c.index]})
	}
}

func (c *Context) parseGroup() {
	groupCtx := &Context{index: c.index + 1, regex: c.regex, Tokens: make([]token, 0, len(c.regex))}
	for c.regex[groupCtx.index] != ')' {
		groupCtx.parseNext()
		groupCtx.index++
	}
	c.index = groupCtx.index
	c.push(TokenGroup{true, groupCtx.Tokens})
}

type charitem struct {
	start  byte
	end    byte
	hasEnd bool
}

func (c *Context) parseCharset() {
	input := c.getRegex()
	items := make([]charitem, 0)
	c.index++
	for input[c.index] != ']' {
		switch input[c.index] {
		case '-':
			items[len(items)-1].end = input[c.index+1]
			items[len(items)-1].hasEnd = true
			c.index++
		default:
			items = append(items, charitem{start: input[c.index]})
			c.index++
		}
	}
	charSet := make(map[byte]struct{})
	for _, item := range items {
		if !item.hasEnd {
			charSet[item.start] = struct{}{}
		} else {
			for i := item.start; i <= item.end; i++ {
				charSet[i] = struct{}{}
			}
		}
	}
	c.push(TokenCharsetLiterals{nil, charSet})
}

func (c *Context) parseAlternation() {
	input := c.getRegex()
	rightSideContext := &Context{index: c.index + 1, regex: c.regex, Tokens: make([]token, 0)}
	for rightSideContext.index < len(input) && input[rightSideContext.index] != ')' {
		rightSideContext.parseNext()
		rightSideContext.index++
	}

	leftToken := TokenGroup{false, c.Tokens} // All existing tokens in current context are on the left side of the alternation
	rightToken := TokenGroup{false, rightSideContext.Tokens}

	c.index = rightSideContext.index - 1       // -1 because the group will increment the index
	c.Tokens = make([]token, 0, cap(c.Tokens)) // Clear old tokens
	c.push(TokenAlternate{leftToken, rightToken})
}

func (c *Context) parseShorthandQuantifier(min int, max int) {
	prevToken := c.Tokens[len(c.Tokens)-1]
	c.Tokens[len(c.Tokens)-1] = TokenQuantifier{ // Replace the last token with a quantifier
		min:   min,
		max:   max,
		token: prevToken,
	}
}

func (c *Context) parseExplicitQuantifier() {
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
	c.parseShorthandQuantifier(min, max)
}
