package window

import (
	"fmt"
	"testing"
)

func Test_consumeRE(t *testing.T) {
	type test struct {
		text, pattern, result string
	}
	tests := []test{
		{"", "", ""},
		{"", "a", ""},
		{"a", "a", "a"},
		{"aab", "a", "a"},
		{"aab", "a+", "aa"},
		{"2 days", `\d+`, "2"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := startParsing(tt.text)
			if r := p.consumeRE(tt.pattern); r != tt.result {
				t.Errorf("unexpected result [%s] when expected [%s]", r, tt.result)
			}
		})
	}
}

func Test_expect(t *testing.T) {
	type test struct {
		text   string
		alts   []string
		result string
	}
	tests := []test{
		{"", []string{""}, ""},
		{"", []string{"a"}, ""},
		{"a", []string{"a"}, "a"},
		{"a", []string{"b"}, ""},
		{"a", []string{"b", "a"}, "a"},
		{"ab", []string{"b", "a"}, "a"},
		{"abc", []string{"a", "ab", "abc"}, "a"}, // picks the first match
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := startParsing(tt.text)
			if r := p.expectAny(tt.alts); r != tt.result {
				t.Errorf("unexpected result [%s] when expected [%s]", r, tt.result)
			}
		})
	}
}

func Test_consumeUntil(t *testing.T) {
	type test struct {
		text              string
		alts              []string
		consumed, matched string
	}
	tests := []test{
		{"", []string{""}, "", ""},
		{"", []string{"a"}, "", ""},
		{"a", []string{"a"}, "", "a"},
		{"abc", []string{"b"}, "a", "b"},
		{"abc", []string{"b", "c"}, "a", "b"}, // picks the first alt
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := startParsing(tt.text)
			if consumed, matched := p.consumeUntil(tt.alts); consumed != tt.consumed || matched != tt.matched {
				t.Errorf("unexpected result [%s, %s] when expected [%s,%s]", consumed, matched, tt.consumed, tt.matched)
			}
		})
	}
}

func Test_eatWs(t *testing.T) {
	type test struct {
		text, result string
	}
	tests := []test{
		{"", ""},
		{" ", ""},
		{"\n", ""},
		{"\t", ""},
		{" a", "a"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := startParsing(tt.text)
			p.eatWs()
			if p.getRemainder() != tt.result {
				t.Errorf("unexpected result [%s] when expected [%s]", p.getRemainder(), tt.result)
			}
		})
	}
}
