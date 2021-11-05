package main

import (
	"flag"
	"io"
	"log"
	"os"
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
	v := newView(w,
		viewWidth(width),
		viewOffset(offset),
		viewTabWidth(tabwidth),
	)
	return v.put(r)
}
