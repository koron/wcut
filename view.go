package main

import (
	"bufio"
	"errors"
	"io"

	"github.com/mattn/go-runewidth"
)

type view struct {
	w  *bufio.Writer
	l  int
	r  int
	tw int

	curr int
}

type viewOption func(*view)

func viewWidth(width int) viewOption {
	return func(v *view) {
		v.r = v.l + width
	}
}

func viewOffset(offset int) viewOption {
	return func(v *view) {
		v.l += offset
		v.r += offset
	}
}

func viewTabWidth(tabwidth int) viewOption {
	return func(v *view) {
		v.tw = tabwidth
	}
}

func newView(w io.Writer, options ...viewOption) *view {
	v := &view{
		w:  bufio.NewWriter(w),
		r:  80,
		tw: 8,
	}
	for _, opt := range options {
		opt(v)
	}
	return v
}

func (v *view) putWSS(n int) error {
	for i := 0; i < n; i++ {
		if v.curr >= v.l && v.curr+1 <= v.r {
			if _, err := v.w.WriteRune(' '); err != nil {
				return err
			}
		}
		v.curr++
	}
	return nil
}

func (v *view) putRune(r rune) error {
	n := runewidth.RuneWidth(r)
	// put broken (partially in-area) rune
	if v.curr < v.l || v.curr+n > v.r {
		return v.putWSS(n)
	}
	// put rune as is.
	if _, err := v.w.WriteRune(r); err != nil {
		return err
	}
	v.curr += n
	return nil
}

func (v *view) put(r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		ru, _, err := br.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		switch ru {
		case '\n':
			if _, err := v.w.WriteRune(ru); err != nil {
				return err
			}
			if err := v.w.Flush(); err != nil {
				return err
			}
			v.curr = 0

		case '\r':

		case '\t':
			err := v.putWSS(v.tw - v.curr%v.tw)
			if err != nil {
				return err
			}

		default:
			err := v.putRune(ru)
			if err != nil {
				return err
			}
		}
	}
	return v.w.Flush()
}
