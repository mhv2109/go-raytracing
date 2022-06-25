package main

import (
	"fmt"
	"io"
)

func WriteColor(w io.Writer, c Color) {
	var r, g, b int

	// scale float values [0, 1] to [0, 255]
	r = int(255.999 * c.X)
	g = int(255.999 * c.Y)
	b = int(255.999 * c.Z)

	fmt.Fprintln(w, r, g, b)
}
