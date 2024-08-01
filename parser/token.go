package parser

import (
	"fmt"
	"strings"
)

type tokenType byte

const (
	groupCaptured tokenType = iota
	groupUncaptured
	alternate
	literal
	quantifier
	charset
)

type token interface {
	getType() tokenType
	getValue() any
	ToString() string
}

type TokenLiteral struct{ rune }

func (t TokenLiteral) getType() tokenType { return literal }
func (t TokenLiteral) getValue() any      { return t.rune }
func (t TokenLiteral) ToString() string   { return fmt.Sprintf("[Literal] %c", t.rune) }

type TokenGroupCaptured struct {
	tokens []token
}

func (t TokenGroupCaptured) getType() tokenType { return groupCaptured }
func (t TokenGroupCaptured) getValue() any      { return t.tokens }
func (t TokenGroupCaptured) ToString() string {
	var sb strings.Builder
	sb.WriteString("[GroupCaptured]")
	for _, token := range t.tokens {
		sb.WriteString("\n    ")
		sb.WriteString(token.ToString())
	}
	return sb.String()
}

type TokenGroupUncaptured struct {
	tokens []token
}

func (t TokenGroupUncaptured) getType() tokenType { return groupUncaptured }
func (t TokenGroupUncaptured) getValue() any      { return t.tokens }
func (t TokenGroupUncaptured) ToString() string {
	var sb strings.Builder
	sb.WriteString("[GroupUncaptured]")
	for _, token := range t.tokens {
		sb.WriteString("\n    ")
		sb.WriteString(token.ToString())
	}
	return sb.String()
}



type TokenCharset struct {
	literals map[byte]struct{}
}

func (t TokenCharset) getType() tokenType { return charset }
func (t TokenCharset) getValue() any      { return t.literals }
func (t TokenCharset) ToString() string {
	var sb strings.Builder
	sb.WriteString("[Literals] ")
	for key := range t.literals {
		sb.WriteByte(key)
	}
	return sb.String()
}

type TokenAlternate struct {
	a token
	b token
}

func (t TokenAlternate) getType() tokenType { return alternate }
func (t TokenAlternate) getValue() any      { return []token{t.a, t.b} }
func (t TokenAlternate) ToString() string {
	var sb strings.Builder
	sb.WriteString("[Alternation]")
	sb.WriteString("\n    ")
	sb.WriteString(t.a.ToString())
	sb.WriteString("\n    ")
	sb.WriteString(t.b.ToString())
	return sb.String()
}
