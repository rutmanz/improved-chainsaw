package parser

type Context struct {
	parent *Context
	regex  *string
	index  int
	Tokens []token
}

func (c *Context) push(token token) {
	c.Tokens = append(c.Tokens, token)
}

func (c *Context) getRegex() string {
	if c.regex == nil {
		return c.parent.getRegex()
	}
	return *c.regex
}



func Parse(regex string) *Context {
	context := &Context{
		index:  0,
		parent: nil,
		regex:  &regex,
		Tokens: make([]token, 0, len(regex)),
	}
	for context.index < len(regex) {
		context.process()
		context.index++
	}
	return context
}
