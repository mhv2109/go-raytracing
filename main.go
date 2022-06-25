package main

import (
	"fmt"
	"math"
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

// rayColor calculates the Color along the Ray. We define objects + colors here,
// and return an object's color if the Ray intersects it. Otherwise, we return
// the background color
func rayColor(r Ray) Color {
	// objects in the scene
	if t := hitSphere(Point3{0, 0, -1}, 0.5, r); t > 0 { // sphere in center of image, with a radius of 0.5
		N := r.At(t).Sub(Point3{0, 0, -1}).Unit()
		return Color{N.X + 1, N.Y + 1, N.Z + 1}.MulS(0.5)
	}

	// background
	var (
		dir = r.Dir.Unit()
		a   = Color{1, 1, 1}       // white
		b   = Color{0.5, 0.7, 1.0} // blue
		t   = 0.5 * (dir.Y + 1.0)
	)
	return a.MulS(1 - t).Add(b.MulS(t)) // (1-t)*white + t*blue
}

// hitSphere checks if r intersects with a spere defined by the center and
// radius. If so, it returns the value of t with (first, camera-facing)
// intersect point of the sphere, or -1.0 otherwise.
func hitSphere(center Vec3, radius float64, r Ray) float64 {
	// A ray intersects the sphere if there exists two solutions for the quadratic
	// equation (P(t) - C) dot (P(t) - C) - r^2 = 0 for all t, where P(t) = A + t*halfb.
	// We can determine this by calulating the descriminant d. This has been
	// simplified using the method in section 6.2.
	var (
		oc = r.Orig.Sub(center) // A - C

		a     = r.Dir.LenSq()
		halfb = oc.Dot(r.Dir)
		c     = oc.LenSq() - radius*radius

		d = halfb*halfb - a*c
	)
	if d < 0 {
		// no real solutions if d < 0
		return -1
	}
	// at least 1 solution if d > 0
	return (-halfb - math.Sqrt(d)) / a
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
			c = rayColor(r)

			writeColor(os.Stdout, c)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
