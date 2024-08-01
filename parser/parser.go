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
	return fmt.Sprintf("[%v] %v", m.index, m.state)
}

func newMatchStep(index int, state *nfaState) *matchStep {
	return &matchStep{state, index}
}

func Match(regex, input string) bool {
	ctx := Parse(regex)
	startNFA := ctx.toNFA()
	currentSteps := []*matchStep{newMatchStep(0, startNFA)}
	for len(currentSteps) != 0 {
		nextSteps := make([]*matchStep, 0)

		for _, step := range currentSteps {
			if step.state.isEnd && step.index == len(input) {
				return true
			}

			for _, transition := range step.state.transitions {
				var passes bool
				if transition.consumes {
					if step.index < len(input) {
						passes = transition.test(input[step.index])
					} else {
						continue
					}
				} else {
					passes = true
				}
				if passes {
					if transition.consumes {
						nextSteps = append(nextSteps, newMatchStep(step.index+1, transition.state))
					} else {
						nextSteps = append(nextSteps, newMatchStep(step.index, transition.state))
					}
				}
			}

		}
		currentSteps = nextSteps
	}

	return false
}
