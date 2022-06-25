package main

import (
	"fmt"
	"os"
)

const (
	// Aspect ratio for Image and Camera
	ratio = 16.0 / 9.0

	// Image
	// This affects the _output_, like the final resolution. Higher values here
	// don't affect the picture _content_, rather the final _quality_.
	imgWidth  = 400
	imgHeight = imgWidth / ratio

	// Camera
	// This affects the _screen_, like focal length and field of view.
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

// RayColor calculates the background color along a Ray.
func RayColor(r Ray) Color {
	var (
		dir = r.Dir.Unit()
		a   = Color{1, 1, 1}       // white
		b   = Color{0.5, 0.7, 1.0} // blue
		t   = 0.5 * (dir.Y + 1.0)
	)
	return a.MulS(1 - t).Add(b.MulS(t)) // (1-t)*white + t*blue
}

func main() {
	fmt.Println("P3")
	fmt.Println(imgWidth, imgHeight)
	fmt.Println("255")

	var (
		u, v float64

		// Ray extrapolates the _sceen_ (see "Camera" above) from the
		// cartesian coordinate of each pixel from the output file (see
		// "Image" above).
		r Ray

		c Color
	)

	// Pan across each pixel of the output image and calculate the color of each.
	for j := imgHeight; j >= 0; j-- {
		fmt.Fprint(os.Stderr, "\rScanlines remaining:", j)
		for i := 0; i < imgWidth; i++ {
			u = float64(i) / (imgWidth - 1)
			v = float64(j) / (imgHeight - 1)
			r = Ray{origin, lowerLeftCorner.Add(horiz.MulS(u), vert.MulS(v), origin.Neg())}
			c = RayColor(r)

			WriteColor(os.Stdout, c)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
