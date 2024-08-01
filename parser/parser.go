package parser

import "fmt"

type Context struct {
	regex  string
	index  int
	Tokens []token
}

func (c *Context) push(token token) {
	c.Tokens = append(c.Tokens, token)
}

func (c *Context) getRegex() string {
	return c.regex
}

func Parse(regex string) *Context {
	context := &Context{
		index:  0,
		regex:  regex,
		Tokens: make([]token, 0, len(regex)),
	}
	for context.index < len(regex) {
		context.parseNext()
		context.index++
	}
	return context
}

type matchStep struct {
	state *nfaState
	index int
}

func (m matchStep) String() string {
	return fmt.Sprintf("%v %v", m.index, m.state)
}

func newMatchStep(index int, state *nfaState) *matchStep {
	return &matchStep{state, index}
}

func Match(regex, input string) bool {
	ctx := Parse(regex)
	startNFA := ctx.toNFA()

	currentSteps := []*matchStep{newMatchStep(0, startNFA)}
	for len(currentSteps) != 0 {
		fmt.Println(currentSteps)
		nextSteps := make([]*matchStep, 0)
		for _, step := range currentSteps {
			if step.state.transitions[input[step.index]] != nil {
				for _, s := range step.state.transitions[input[step.index]] {
					if s.isEnd {
						return true
					}
					if step.index+1 < len(input) {
						nextSteps = append(nextSteps, newMatchStep(step.index+1, s))
					}
				}
			}
			if step.state.transitions[empty] != nil {
				for _, s := range step.state.transitions[empty] {
					if s.isEnd {
						return true
					}
					nextSteps = append(nextSteps, newMatchStep(step.index, s))
				}
			}
		}
		currentSteps = nextSteps
	}

	return false
}
