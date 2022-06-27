package main

import (
	"fmt"
	"io"
)

func writeColor(w io.Writer, c Color, samples int) {
	// Divide the color by the number of samples and scale float values [0, 1]
	// to [0, 255]
	var (
		scale = 1.0 / float64(samples)
		r     = int(256 * clamp(c.X*scale, 0.0, 0.999))
		g     = int(256 * clamp(c.Y*scale, 0.0, 0.999))
		b     = int(256 * clamp(c.Z*scale, 0.0, 0.999))
	)
	fmt.Fprintln(w, r, g, b)
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}
