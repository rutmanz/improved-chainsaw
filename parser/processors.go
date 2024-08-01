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

func createCharsetLiterals(chars ...byte) TokenCharsetLiterals {
	charSet := make(map[byte]struct{}, len(chars))
	for _, char := range chars {
		charSet[char] = struct{}{}
	}
	return TokenCharsetLiterals{&TokenCharset{}, charSet}
}

var (
	// Charsets
	charsetDigits        interfaceCharset = TokenCharsetRange{&TokenCharset{}, '0', '9'}
	charsetNotDigits     interfaceCharset = TokenCharsetNot{&TokenCharset{}, charsetDigits}
	charsetWhitespace    interfaceCharset = createCharsetLiterals('\f', '\n', '\r', '\t', '\v', '\u0020', '\u00A0')
	charsetNotWhitespace interfaceCharset = TokenCharsetNot{&TokenCharset{}, charsetWhitespace}
	charsetWord          interfaceCharset = TokenCharsetCompound{&TokenCharset{}, []interfaceCharset{charsetDigits, TokenCharsetRange{&TokenCharset{}, 'a', 'z'}, TokenCharsetRange{&TokenCharset{}, 'A', 'Z'}, TokenCharsetLiterals{&TokenCharset{}, map[byte]struct{}{'_': {}}}}}
	charsetNotWord       interfaceCharset = TokenCharsetNot{&TokenCharset{}, charsetWord}
	charsetAny           interfaceCharset = TokenCharsetNot{&TokenCharset{}, createCharsetLiterals('\n', '\r')}
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
	items := make([]interfaceCharset, 0)
	literals := make([]byte, 0)
	c.index++
	for input[c.index] != ']' {
		switch input[c.index] {
		case '-':
			start := literals[len(literals)-1]
			literals = literals[:len(literals)-1]
			items = append(items,TokenCharsetRange{&TokenCharset{}, start, input[c.index+1]})
			c.index++
		default:
			literals = append(literals, input[c.index])
			c.index++
		}
	}
	if (len(literals) > 0) {
		items = append(items, createCharsetLiterals(literals...))
	}
	c.push(TokenCharsetCompound{&TokenCharset{}, items})
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
	origIndex := c.index
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
			panic("Invalid quantifier " + input[origIndex:c.index+1])
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
