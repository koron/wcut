package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWcut(t *testing.T) {
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
			t.Errorf("wcut: width=%q, offset=%d, tabwidth=%d, %d) = %q, want %q", test.input, test.width, test.offset, test.tabwidth, buf.String(), test.output)
		}
	}
}
