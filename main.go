package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

const (
	// Image
	// This affects the _output_, like the final resolution. Higher values here
	// don't affect the picture _content_, rather the final _quality_.
	imgWidth  = 400
	imgHeight = imgWidth / ratio
	samples   = 100
)

var (
	// objects in the scene
	world = NewHittables(
		Sphere{Point3{0, 0, -1}, 0.5},      // sphere in center of image, with a radius of 0.5
		Sphere{Point3{0, -100.5, -1}, 100}, // ground
	)
)

// rayColor calculates the Color along the Ray. We define objects + colors here,
// and return an object's color if the Ray intersects it. Otherwise, we return
// the background color
func rayColor(r Ray) Color {
	// objects in the scene
	if hr := world.Hit(r, 0, math.MaxFloat64); hr != nil {
		return hr.N.Add(Color{1, 1, 1}).MulS(0.5)
	}

	// if no object hit, render background
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

		cam = NewCamera(Point3{0, 0, 0})

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
			pixel := Color{0, 0, 0}
			for s := 0; s < samples; s++ {
				u = (float64(i) + rand.Float64()) / (imgWidth - 1)
				v = (float64(j) + rand.Float64()) / (imgHeight - 1)
				r = cam.Ray(u, v)
				c = rayColor(r)
				pixel = pixel.Add(c)
			}
			writeColor(os.Stdout, pixel, samples)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
