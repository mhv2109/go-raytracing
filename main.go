package main

import (
	"fmt"
	"os"
)

const (
	imgWidth  = 256
	imgHeight = 256
)

func main() {
	fmt.Println("P3")
	fmt.Println(imgWidth, imgHeight)
	fmt.Println("255")

	var (
		r, g, b    float32
		ir, ig, ib int
	)

	for j := imgHeight; j >= 0; j-- {
		fmt.Fprint(os.Stderr, "\rScanlines remaining:", j)
		for i := 0; i < imgWidth; i++ {
			r = float32(i) / (imgWidth - 1)
			g = float32(j) / (imgHeight - 1)
			b = 0.25

			ir = int(255.999 * r)
			ig = int(255.999 * g)
			ib = int(255.99 * b)

			fmt.Println(ir, ig, ib)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
