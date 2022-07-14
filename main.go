package main

import (
	"flag"
	"fmt"
	"io"
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

// diffustionMaterial allows us to select the diffusion function at runtime
func diffusionMaterial() DiffusionOpt {
	if simpleDiff {
		return WithDiffusionType(SimpleDiffusion)
	}
	return WithDiffusionType(Lambertian)
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

func writeColor(w io.Writer, c Color, samples int) {
	// Divide the color by the number of samples and scale float values [0, 1]
	// to [0, 255]
	var (
		scale = 1.0 / float64(samples)
		r     = int(256 * clamp(math.Sqrt(c.X*scale), 0.0, 0.999))
		g     = int(256 * clamp(math.Sqrt(c.Y*scale), 0.0, 0.999))
		b     = int(256 * clamp(math.Sqrt(c.Z*scale), 0.0, 0.999))
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

func main() {
	flag.Parse()

	// build world

	var (
		// materials + surfaces
		ground = NewDiffusion(Color{0.8, 0.8, 0}, diffusionMaterial())
		center = NewDiffusion(Color{0.1, 0.2, 0.5}, diffusionMaterial())
		left   = NewDielectric(Color{1.0, 1.0, 1.0}, IndexOfRefraction(1.5))
		right  = NewMetal(Color{0.8, 0.6, 0.2}, Fuzz(0.0))

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

		lookfrom  = Point3{3, 3, 2}
		lookat    = Point3{0, 0, -1}
		vup       = Vec3{0, 1, 0}
		vfov      = 20.0
		aperture  = 2.0
		focusDist = lookfrom.Sub(lookat).Len()
		cam       = NewCamera(lookfrom, lookat, vup, vfov, aperture, focusDist)

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
