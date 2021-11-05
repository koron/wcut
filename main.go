package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"github.com/mattn/go-runewidth"
)

var (
	width    int
	offset   int
	tabwidth int
)

func main() {
	flag.IntVar(&width, "w", 80, `width to cut`)
	flag.IntVar(&offset, "o", 0, `horizontal offset`)
	flag.IntVar(&tabwidth, "t", 8, `tab width`)
	flag.Parse()
	if tabwidth <= 0 {
		log.Fatalf("tab width should be larger than zero")
	}

	var r io.Reader = os.Stdin
	if flag.NArg() > 0 {
		// FIXME: better multiple files reader.
		rr := make([]io.Reader, 0, flag.NArg())
		for _, name := range flag.Args() {
			f, err := os.Open(name)
			if err != nil {
				log.Fatalf("failed to open file: %s", err)
			}
			defer f.Close()
			rr = append(rr, f)
		}
		r = io.MultiReader(rr...)
	}

	err := wcut(os.Stdout, r)
	if err != nil {
		log.Fatalf("wcut failed: %s", err)
	}
}

func wcut(w io.Writer, r io.Reader) error {
	bw := bufio.NewWriter(w)
	br := bufio.NewReader(r)
	left := offset
	right := offset + width
	curr := 0
	for {
		ru, _, err := br.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil
		}

		switch ru {
		case '\n':
			if _, err := bw.WriteRune(ru); err != nil {
				return err
			}
			if err := bw.Flush(); err != nil {
				return err
			}
			curr = 0

		case '\r':

		case '\t':
			n := tabwidth - curr%tabwidth
			for i := 0; i < n; i++ {
				if curr >= left && curr+1 <= right {
					if _, err := bw.WriteRune(' '); err != nil {
						return err
					}
				}
				curr++
			}

		default:
			n := runewidth.RuneWidth(ru)
			if curr >= left && curr+n <= right {
				if _, err := bw.WriteRune(ru); err != nil {
					return err
				}
				curr += n
				break
			}
			// write broken (partially in-area) rune
			for i := 0; i < n; i++ {
				if curr >= left && curr+1 <= right {
					if _, err := bw.WriteRune(' '); err != nil {
						return err
					}
				}
				curr++
			}
		}
	}
	return bw.Flush()
}
