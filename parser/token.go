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
	ToString(padding string) string
}

type TokenLiteral struct{ rune }

func (t TokenLiteral) getType() tokenType { return literal }
func (t TokenLiteral) ToString(padding string) string {
	return fmt.Sprintf("%v[Literal] %c", padding, t.rune)
}

type TokenGroupCaptured struct {
	tokens []token
}

func (t TokenGroupCaptured) getType() tokenType { return groupCaptured }
func (t TokenGroupCaptured) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + "[GroupCaptured]")
	for _, token := range t.tokens {
		sb.WriteString("\n" + token.ToString(padding+"    "))
	}
	return sb.String()
}

type TokenGroupUncaptured struct {
	tokens []token
}

func (t TokenGroupUncaptured) getType() tokenType { return groupUncaptured }
func (t TokenGroupUncaptured) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + "[GroupUncaptured]")
	for _, token := range t.tokens {
		sb.WriteString("\n" + token.ToString(padding+"    "))
	}
	return sb.String()
}

type interfaceCharset interface {
	token
	Test(byte) bool
}
type TokenCharset struct {
}

func (t TokenCharset) getType() tokenType { return charset }


type TokenCharsetAny struct{*TokenCharset}
func (t TokenCharsetAny) Test(b byte) bool               { return true }
func (t TokenCharsetAny) ToString(padding string) string { return padding + "[Charset] Wildcard" }

type TokenCharsetRange struct{
	*TokenCharset
	min byte
	max byte
}
func (t TokenCharsetRange) Test(b byte) bool               { return b >= t.min && b <= t.max }
func (t TokenCharsetRange) ToString(padding string) string { return padding + fmt.Sprintf("[Charset] %c - %c", t.min, t.max) }

type TokenCharsetNot struct {
	*TokenCharset
	tester interfaceCharset
}
func (t TokenCharsetNot) Test(b byte) bool               { return !t.tester.Test(b) }
func (t TokenCharsetNot) ToString(padding string) string { 

	return padding + "[Charset] Not\n" + t.tester.ToString(padding+"    ")
}

type TokenCharsetLiterals struct {
	*TokenCharset
	literals map[byte]struct{}
}

func (t TokenCharsetLiterals) Test(b byte) bool {
	_, ok := t.literals[b]
	return ok
}
func (t TokenCharsetLiterals) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + "[Charset]")
	for b := range t.literals {
		sb.WriteString(fmt.Sprintf("\n%s    %c", padding, b))
	}
	return sb.String()
}

type TokenAlternate struct {
	a token
	b token
}

func (t TokenAlternate) getType() tokenType { return alternate }
func (t TokenAlternate) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + "[Alternation]")
	sb.WriteString("\n" + t.a.ToString(padding+"    "))
	sb.WriteString("\n" + t.b.ToString(padding+"    "))
	return sb.String()
}

type TokenQuantifier struct {
	min   int
	max   int
	token token
}

func (t TokenQuantifier) getType() tokenType { return quantifier }
func (t TokenQuantifier) ToString(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding + "[Quantifier ")
	if t.max != -1 {
		sb.WriteString(fmt.Sprintf("%d-%d", t.min, t.max))
	} else {
		sb.WriteString(fmt.Sprintf("%d+", t.min))
	}
	sb.WriteString("]")
	sb.WriteString("\n" + t.token.ToString(padding+"    "))
	return sb.String()
}
