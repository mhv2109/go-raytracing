package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

const (
	// Image
	// This affects the _output_, like the final resolution. Higher values here
	// don't affect the picture _content_, rather the final _quality_.
	imgWidth  = 2560
	imgHeight = int(imgWidth / ratio)
	samples   = 500
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

func randomScene() Hittables {
	var (
		world  = NewHittables()
		ground = NewDiffusion(Color{0.8, 0.8, 0}, diffusionMaterial())
	)

	world.Add(Sphere{Point3{0, -1000, 0}, 1000, ground})

	p := Point3{4, 0.2, 0}
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Point3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}
			if center.Sub(p).Len() > 0.9 {
				var (
					c Color
					m Material
				)

				switch choose := rand.Float64(); {
				case choose < 0.8:
					// diffuse
					c = RandomVec3(0, 1).Mul(RandomVec3(0, 1))
					m = NewDiffusion(c, diffusionMaterial())
				case choose < 0.95:
					// metal
					c = RandomVec3(0.5, 1)
					fuzz := rand.Float64() * 0.5
					m = NewMetal(c, Fuzz(fuzz))
				default:
					// glass
					m = NewDielectric(Color{1, 1, 1}, IndexOfRefraction(1.5))
				}

				world.Add(Sphere{center, 0.2, m})
			}
		}
	}

	material1 := NewDielectric(Color{1, 1, 1}, IndexOfRefraction(1.5))
	world.Add(Sphere{Point3{0, 1, 0}, 1, material1})

	material2 := NewDiffusion(Color{0.4, 0.2, 0.1}, diffusionMaterial())
	world.Add(Sphere{Point3{-4, 1, 0}, 1, material2})

	material3 := NewMetal(Color{0.7, 0.6, 0.5}, Fuzz(0))
	world.Add(Sphere{Point3{4, 1, 0}, 1, material3})

	return world
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

func process(cam Camera, world Hittables) {
	// Pan across each pixel of the output image and calculate the color of each.
	var (
		wg      sync.WaitGroup
		results = make(chan chan Color, 2*runtime.NumCPU())
	)

	wg.Add(1)
	go func() {
		for j := imgHeight; j >= 0; j-- {
			for i := 0; i < imgWidth; i++ {
				// calculate each ray concurrently
				ch := make(chan Color, 1)
				results <- ch

				go func(j, i int, ch chan Color) {
					var (
						u, v  float64
						pixel = Color{0, 0, 0}
						r     Ray
						c     Color
					)

					for s := 0; s < samples; s++ {
						u = (float64(i) + rand.Float64()) / (imgWidth - 1)
						v = (float64(j) + rand.Float64()) / (float64(imgHeight) - 1)
						r = cam.Ray(u, v)
						c = rayColor(r, world)
						pixel = pixel.Add(c)
					}

					ch <- pixel
					close(ch)
				}(j, i, ch)
			}
		}
		close(results)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for ch := range results {
			pixel := <-ch
			writeColor(os.Stdout, pixel, samples)
		}
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	flag.Parse()

	// build world
	world := randomScene()

	// output image

	fmt.Println("P3")
	fmt.Println(imgWidth, imgHeight)
	fmt.Println("255")

	var (
		lookfrom  = Point3{13, 2, 3}
		lookat    = Point3{0, 0, -0}
		vup       = Vec3{0, 1, 0}
		vfov      = 20.0
		aperture  = 0.1
		focusDist = 10.0
		cam       = NewCamera(lookfrom, lookat, vup, vfov, aperture, focusDist)
	)

	process(cam, world)
}
