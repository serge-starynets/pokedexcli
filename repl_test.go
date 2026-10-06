package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  l o L  ",
			expected: []string{"l", "o", "l"},
		},
		{
			input:    "  hello  world SuPer  DUPER ",
			expected: []string{"hello", "world", "super", "duper"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Lengths are different")
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("word is different")
				t.Fail()
			}
		}
	}
}
