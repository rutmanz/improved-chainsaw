package parser

import (
	"fmt"
)

type nfaTransitions = map[byte][]*nfaState
type nfaState struct {
	isStart     bool
	isEnd       bool
	transitions nfaTransitions
}

func (n *nfaState) ToString(padding string) string {
	return fmt.Sprintf("%v", n.transitions)
}

func newNfaState() *nfaState {
	return &nfaState{transitions: make(nfaTransitions)}
}

func (n *nfaState) addTransition(char byte, states ...*nfaState) {
	if n.transitions[char] == nil {
		n.transitions[char] = make([]*nfaState, 0)
	}
	n.transitions[char] = append(n.transitions[char], states...)
}
func (n *nfaState) setTransitions(char byte, states ...*nfaState) {
	n.transitions[char] = states
}

const empty byte = 0

func (c *Context) toNFA() *nfaState {
	currStart, currEnd := toStates(c.Tokens[0])
	for i := 1; i < len(c.Tokens); i++ {
		nextStart, nextEnd := toStates(c.Tokens[i])
		currEnd.addTransition(empty, nextStart)
		currEnd = nextEnd
	}

	start := newNfaState()
	start.addTransition(empty, currStart)
	start.isStart = true

	end := newNfaState()
	end.isEnd = true

	currEnd.addTransition(empty, end)

	return start
}

func toStates(t token) (start *nfaState, end *nfaState) {
	start = newNfaState()
	end = newNfaState()
	switch t.getType() {
	case groupCaptured, groupUncaptured:
		t := t.(TokenGroup)
		start, end = toStates(t.tokens[0])
		for i := 1; i < len(t.tokens); i++ {
			nextStart, nextEnd := toStates(t.tokens[i])
			end.addTransition(empty, nextStart)
			end = nextEnd
		}
	case alternate:
		t := t.(TokenAlternate)
		startA, endA := toStates(t.a)
		startB, endB := toStates(t.b)
		start.setTransitions(empty, startA, startB)
		endA.setTransitions(empty, end)
		endB.setTransitions(empty, end)
	case literal:
		t := t.(TokenLiteral)
		start.setTransitions(t.byte, end)
	case charset:
		// t := t.(interfaceCharset)
	case quantifier:
		// t := t.(TokenQuantifier)

	default:
		panic("Unknown token type")
	}
	return
}
