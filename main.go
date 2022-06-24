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
		r, g, b float64
		c       Color
	)

	for j := imgHeight; j >= 0; j-- {
		fmt.Fprint(os.Stderr, "\rScanlines remaining:", j)
		for i := 0; i < imgWidth; i++ {
			r = float64(i) / (imgWidth - 1)
			g = float64(j) / (imgHeight - 1)
			b = 0.25
			c = Color{r, g, b}

			WriteColor(os.Stdout, c)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
