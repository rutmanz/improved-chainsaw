package parser_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/rutmanz/improved-chainsaw/parser"
)

var tests = []struct {
	regex string
	input string
}{
	{"(a|b)c", "bc"},
	{"\\wc", "bc"},
	{"(a|b)c", "dsac"},
	{"(a|b)+c", "aac"},
}

func TestGeneric(t *testing.T) {
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v,%v", test.regex, test.input), func(t *testing.T) {
			expected := regexp.MustCompile("^" + test.regex + "$").MatchString(test.input)
			actual := parser.Match(test.regex, test.input)
			if expected != actual {
				t.Fatalf("Expected %t, got %t | %v", expected, actual, test)
			} else {
				t.Logf("Pattern: %v | Value: %v", test.regex, test.input)
			}
		})

	}
}
