package parser

import (
	"fmt"
	"strings"
)

type nfaTransition struct {
	test     func(byte) (passes bool)
	consumes bool
	state    *nfaState
	label    string
}

func (n *nfaTransition) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + fmt.Sprintf("%-5v", "{"+n.label+"}"))
	sb.WriteString(n.state.ToString(padding))
	return sb.String()
}

type nfaState struct {
	isStart     bool
	isEnd       bool
	transitions []*nfaTransition
}

func (n *nfaState) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + fmt.Sprintf("[State] Start=%t End=%t", n.isStart, n.isEnd))
	if len(n.transitions) == 1 {
		sb.WriteString("\n" + n.transitions[0].ToString(padding+"  "))
	} else {
		for _, transition := range n.transitions {
			sb.WriteString("\n" + transition.ToString(padding+"  "))
		}
	}
	return sb.String()
}

func (n *nfaState) String() string {
	return n.ToString("")
}

func newNfaState() *nfaState {
	return &nfaState{transitions: make([]*nfaTransition, 0)}
}

func (n *nfaState) addEpsilonTransition(states ...*nfaState) {
	for _, state := range states {
		n.transitions = append(n.transitions, &nfaTransition{func(b byte) bool { return true }, false, state, "*"})
	}
}
func (n *nfaState) addTransition(label string, test func(byte) bool, consumes bool, states ...*nfaState) {
	for _, state := range states {
		n.transitions = append(n.transitions, &nfaTransition{test, consumes, state, label})
	}
}

func (c *Context) toNFA() *nfaState {
	currStart, currEnd := toStates(c.Tokens[0])
	for i := 1; i < len(c.Tokens); i++ {
		nextStart, nextEnd := toStates(c.Tokens[i])
		currEnd.addEpsilonTransition(nextStart)
		currEnd = nextEnd
	}

	start := newNfaState()
	start.addEpsilonTransition(currStart)
	start.isStart = true

	end := newNfaState()
	end.isEnd = true

	currEnd.addEpsilonTransition(end)

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
			end.addEpsilonTransition(nextStart)
			end = nextEnd
		}
	case alternate:
		t := t.(TokenAlternate)
		startA, endA := toStates(t.a)
		startB, endB := toStates(t.b)
		start.addEpsilonTransition(startA, startB)
		endA.addEpsilonTransition(end)
		endB.addEpsilonTransition(end)
	case literal:
		t := t.(TokenLiteral)
		start.addTransition(fmt.Sprintf("<%c>", t.byte), func(b byte) bool { return b == t.byte }, true, end)
	case charset:
		t := t.(interfaceCharset)
		start.addTransition(fmt.Sprintf("<%v>", t.ToString("")), t.Test, true, end)
	case quantifier:
		t := t.(TokenQuantifier)
		if t.min == 0 {
			start.addEpsilonTransition(end)
		}
		instances := t.min
		if t.max == -1 {
			if t.min == 0 {
				instances = 1
			}
		} else {
			instances = t.max
		}
		tStart, tEnd := toStates(t.token)
		start.addEpsilonTransition(tStart)
		for i := 2; i <= instances; i++ {
			iterStart, iterEnd := toStates(t.token)
			tEnd.addEpsilonTransition(iterStart)
			tEnd = iterEnd
			tStart = iterStart
			if i > t.min {
				iterStart.addEpsilonTransition(end)
			}
		}
		tEnd.addEpsilonTransition(end)
		if t.max == -1 {
			end.addEpsilonTransition(tStart)
		}

	default:
		panic("Unknown token type")
	}
	return
}
