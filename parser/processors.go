package parser

func (c *Context) process() {
	input := c.getRegex()
	switch input[c.index] {
	case '(':
		groupContext := &Context{index: c.index + 1, parent: c, Tokens: make([]token, 0, len(input))}
		c.push(groupContext.processGroup())
	case '[':
		c.push(c.processCharset())
	case '|':
		c.push(c.processAlternation())
	case '*':
	case '?':
	case '+':
		break
	case '{':
		break
	default:
		c.push(TokenLiteral{rune(input[c.index])})
	}
}

func (c *Context) processGroup() token {
	input := c.getRegex()
	for input[c.index] != ')' {
		c.process()
		c.index++
	}
	c.parent.index = c.index
	return TokenGroupCaptured{c.Tokens}
}

type charitem struct {
	start byte
	end   byte
}

func (c *Context) processCharset() token {
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
	return TokenCharset{charSet}
}

func (c *Context) processAlternation() token {
	input := c.getRegex()
	rightSideContext := &Context{index: c.index + 1, parent: c, Tokens: make([]token, 0, len(input)-c.index)}
	for rightSideContext.index < len(input) && input[rightSideContext.index] != '(' {
		rightSideContext.process()
		rightSideContext.index++
	}

	leftToken := TokenGroupUncaptured{c.Tokens} // All existing tokens in current context are on the left side of the alternation
	rightToken := TokenGroupUncaptured{rightSideContext.Tokens}

	c.index = rightSideContext.index
	c.Tokens = make([]token, 0, cap(c.Tokens)) // Clear old tokens
	return TokenAlternate{leftToken, rightToken}
}
