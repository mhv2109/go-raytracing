package main

import (
	"fmt"
	"os"
)

const (
	ratio = 16.0 / 9.0

	// Image
	imgWidth  = 400
	imgHeight = imgWidth / ratio

	// Camera
	viewHeight = 2.0
	viewWidth  = ratio * viewHeight
	focalLen   = 1.0
)

var (
	// Camera
	origin          = Point3{0, 0, 0}
	horiz           = Point3{viewWidth, 0, 0}
	vert            = Point3{0, viewHeight, 0}
	lowerLeftCorner = origin.Sub(horiz.DivS(2), vert.DivS(2), Point3{0, 0, focalLen})
)

// RayColor calculates the background color along a Ray
func RayColor(r Ray) Color {
	var (
		dir = r.Dir.Unit()
		a   = Color{1, 1, 1}
		b   = Color{0.5, 0.7, 1.0}
		t   = 0.5 * (dir.Y + 1.0)
	)
	return a.MulS(1 - t).Add(b.MulS(t))
}

func main() {
	fmt.Println("P3")
	fmt.Println(imgWidth, imgHeight)
	fmt.Println("255")

	for j := imgHeight; j >= 0; j-- {
		fmt.Fprint(os.Stderr, "\rScanlines remaining:", j)
		for i := 0; i < imgWidth; i++ {
			var (
				u = float64(i) / (imgWidth - 1)
				v = float64(j) / (imgHeight - 1)
				r = Ray{origin, lowerLeftCorner.Add(horiz.MulS(u), vert.MulS(v), origin.Neg())}
				c = RayColor(r)
			)
			WriteColor(os.Stdout, c)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
