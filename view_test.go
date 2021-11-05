package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestView(t *testing.T) {
	tests := []struct {
		input    string
		output   string
		width    int
		offset   int
		tabwidth int
	}{
		{
			input:    `ほげーーーーーーーーーーーー`,
			output:   `ほげーーー`,
			width:    10,
			offset:   0,
			tabwidth: 4,
		},
		{
			input:    `ほげーーーーーーーーーーーー`,
			output:   `げーーーー`,
			width:    10,
			offset:   2,
			tabwidth: 4,
		},
		{
			input:    `ほげーーーーーーーーーーーー`,
			output:   `ほげ `,
			width:    5,
			offset:   0,
			tabwidth: 4,
		},
		{
			input:    `ほげーーーーーーーーーーーー`,
			output:   `ほげ`,
			width:    4,
			offset:   0,
			tabwidth: 4,
		},
		{
			input:    `ほげーーーーーーーーーーーー`,
			output:   ` げーーー `,
			width:    10,
			offset:   1,
			tabwidth: 4,
		},

		// including TAB ('\t')
		{
			input:    "abc\tdef",
			output:   "abc def",
			width:    10,
			offset:   0,
			tabwidth: 4,
		},
		{
			input:    "abc\tdef",
			output:   "abc     de",
			width:    10,
			offset:   0,
			tabwidth: 8,
		},
		{
			input:    "abc\tdef",
			output:   "bc     def",
			width:    10,
			offset:   1,
			tabwidth: 8,
		},
		{
			input:    "abc\tdef",
			output:   "c     def",
			width:    10,
			offset:   2,
			tabwidth: 8,
		},
		{
			input:    "abc\tdef",
			output:   "     def",
			width:    10,
			offset:   3,
			tabwidth: 8,
		},
		{
			input:    "abc\tdef",
			output:   "    def",
			width:    10,
			offset:   4,
			tabwidth: 8,
		},

		// including line feed ('\n')
		{
			input:    "123456789A123456789B123456789C123456789D\nfoo bar baz qux quux",
			output:   "12345\nfoo b",
			width:    5,
			offset:   0,
			tabwidth: 8,
		},
	}
	for _, test := range tests {
		var buf bytes.Buffer
		err := newView(&buf,
			viewWidth(test.width),
			viewOffset(test.offset),
			viewTabWidth(test.tabwidth),
		).put(strings.NewReader(test.input))
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != test.output {
			t.Errorf("view.put failed: input=%q, width=%d, offset=%d, tabwidth=%d: want %q but got %q", test.input, test.width, test.offset, test.tabwidth, test.output, buf.String())
		}
	}
}
