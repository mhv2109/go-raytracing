package main

import (
	"flag"
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
	imgHeight = int(imgWidth / ratio)
	samples   = 100
)

var (
	// cmdline args
	simpleDiff bool
)

func init() {
	flag.BoolVar(&simpleDiff, "simple", false, "use simple diffusion calculation")
}

// rayColor calculates the Color along the Ray. We define objects + colors here,
// and return an object's color if the Ray intersects it. Otherwise, we return
// the background color
func rayColor(r Ray, world Hittables) Color {
	var (
		mult = Vec3{1, 1, 1}
		hr   *HitRecord
		n    = 0
	)

LOOP:
	if n > 50 {
		return Color{0, 0, 0}
	}

	hr = nil
	if hr = world.Hit(r, 1e-3, math.MaxFloat64); hr == nil {
		// if no object hit, render background
		var (
			dir = r.Dir.Unit()
			a   = Color{1, 1, 1}       // white
			b   = Color{0.5, 0.7, 1.0} // blue
			t   = 0.5 * (dir.Y + 1.0)
		)
		return a.MulS(1 - t).Add(b.MulS(t)).Mul(mult) // (1-t)*white + t*blue
	}

	// objects in the scene
	att, scatt := hr.M.Scatter(r, *hr)
	if att == nil || scatt == nil {
		return Color{0, 0, 0}
	}
	r = *scatt
	mult = mult.Mul(*att)

	n++
	goto LOOP // recursive version causes stack overflow
}

// diffustionMaterial allows us to select the diffusion function at runtime
func diffusionMaterial() MaterialType {
	if simpleDiff {
		return SimpleDiffusion
	}
	return Lambertian
}

func main() {
	flag.Parse()

	// build world

	var (
		// materials + surfaces
		ground = NewMaterial(Color{0.8, 0.8, 0}, diffusionMaterial())
		center = NewMaterial(Color{0.1, 0.2, 0.5}, diffusionMaterial())
		left   = NewMaterial(Color{1.0, 1.0, 1.0}, Dielectric, DielectricIndexOfRefraction(1.5))
		right  = NewMaterial(Color{0.8, 0.6, 0.2}, Metal, MetalFuzz(0.0))

		// objects in the scene
		world = NewHittables(
			Sphere{Point3{0, -100.5, -1}, 100, ground}, // ground
			Sphere{Point3{0, 0, -1}, 0.5, center},      // sphere in center of image, with a radius of 0.5
			Sphere{Point3{-1, 0, -1}, 0.5, left},
			Sphere{Point3{-1, 0, -1}, -0.4, left}, // hollow sphere. See 10.5
			Sphere{Point3{1, 0, -1}, 0.5, right},
		)
	)

	// output image

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
				v = (float64(j) + rand.Float64()) / (float64(imgHeight) - 1)
				r = cam.Ray(u, v)
				c = rayColor(r, world)
				pixel = pixel.Add(c)
			}
			writeColor(os.Stdout, pixel, samples)
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}
