package parser_test

import (
	"regexp"
	"testing"

	"github.com/rutmanz/improved-chainsaw/parser"
)

var tests = []struct {
	regex  string
	inputs []string
}{
	// Quantifiers
	{"a+", []string{"", "a", "aa", "aaaaaaa", "aaaaaba"}},
	{"a*", []string{"", "a", "aa", "aaaaaaa", "aaaaaba"}},
	{"a?", []string{"", "a", "aa", "aaaaaaa", "aaaaaba"}},
	{"a{2,3}", []string{"", "a", "aa", "aaa", "aaaaa"}},
	{"a{0,2}", []string{"", "a", "aa", "aaa"}},
	{"a{2,}", []string{"a", "aa", "aaaaaaa"}},
	{"a{2}", []string{"a", "aa", "aaa"}},
	{"(a|b)+", []string{"", "aaaaa", "ababbbab", "baba", "b", "aaaaabaaaa"}},
	{"c(a|b)+", []string{"", "c", "caaaaa", "acbabbbab", "cbaba", "cb", "caaaaabaaaa"}},
	{"c(a|b)?", []string{"", "c", "ca", "cb", "cab"}},

	// Charsets
	{"[a-y]", []string{" ", "a", "y", "z", "A", "Y", "Z"}},
	{"[A-y]", []string{" ", "a", "y", "z", "A", "Y", "Z"}},
	{"[A-Y]", []string{" ", "a", "y", "z", "A", "Y", "Z"}},
	{"[1ab]", []string{"1", "a", "b", "1ab", "2", "2b"}},
	{"\\w+", []string{"", "1", "a", "aa", "aaaaaaa", "aa1aa2aba", "*", ":"}},
	{"\\W+", []string{"", "1", "a", "aa", "aaaaaaa", "aa1aa2aba", "*", ":"}},
	{"\\d+", []string{"a", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}},
	{"\\D+", []string{"a", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}},
	{".", []string{".", "a", "1", "_", ":", " ", "\n", "\t"}},
	{"\\.", []string{".", "a", "1", "_", ":", " ", "\n", "\t"}},

	{"a|b|c", []string{"a", "b", "c", "ab"}},
}

func TestGeneric(t *testing.T) {
	for _, test := range tests {
		t.Run(test.regex, func(t *testing.T) {
			compiled, err := regexp.Compile("^(?:" + test.regex + ")$")
			if err != nil {
				t.Errorf("Failed to compile regex: %v", err)
			}
			for _, input := range test.inputs {
				t.Run(input, func(t *testing.T) {
					expected := compiled.MatchString(input)
					actual := parser.Match(test.regex, input)
					if expected != actual {
						t.Errorf("Expected %t, got %t | '%v' |", expected, actual, input)
					} else {
						t.Logf("Pattern: %v | Value: %v", test.regex, input)
					}
				})
			}
		})
	}
}
