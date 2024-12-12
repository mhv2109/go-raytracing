package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

var (
	// cmdline args
	imgWidth   int
	imgHeight  int
	samples    int
	jobs       int
	simpleDiff bool
	cpuprofile string
)

func init() {
	flag.IntVar(&imgWidth, "width", 2560, "output image width")
	flag.IntVar(&imgHeight, "height", 1440, "output image height")
	flag.IntVar(&samples, "samples", 500, "number of samples per pixel")
	flag.IntVar(&jobs, "jobs", 2*runtime.NumCPU(), "number of jobs for rendering")
	flag.BoolVar(&simpleDiff, "simple", false, "use simple diffusion calculation")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "create a CPU profile and save to file")
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
		mult  = Vec3{1, 1, 1}
		hr    HitRecord
		n     = 0
		att   Color
		scatt Ray
	)

LOOP:
	if n > 50 {
		return Color{0, 0, 0}
	}

	if !world.Hit(r, 1e-3, math.MaxFloat64, &hr) {
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
	if !hr.M.Scatter(r, hr, &att, &scatt) {
		return Color{0, 0, 0}
	}
	r = scatt
	mult = mult.Mul(att)

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

func process(cam Camera, world Hittables, w, h, ns int, writer io.Writer) {
	fmt.Fprintln(writer, "P3")
	fmt.Fprintln(writer, w, h)
	fmt.Fprintln(writer, "255")

	// Pan across each pixel of the output image and calculate the color of each.
	var (
		wg      sync.WaitGroup
		results = make(chan chan Color, jobs)
		rands   = make([]float64, 2*ns)
	)

	// calculate random values up-front
	for i := range rands {
		rands[i] = <-RandomCh
	}

	wg.Add(1)
	go func() {
		for j := h; j >= 0; j-- {
			for i := 0; i < w; i++ {
				// calculate each ray concurrently
				ch := make(chan Color, 1)
				results <- ch

				go func(j, i int) {
					var (
						u, v  float64
						pixel = Color{0, 0, 0}
						r     Ray
						c     Color
					)

					for s := 0; s < 2*ns; s += 2 {
						u = (float64(i) + rands[s]) / (float64(w) - 1)
						v = (float64(j) + rands[s+1]) / (float64(h) - 1)
						r = cam.Ray(u, v)
						c = rayColor(r, world)
						pixel = pixel.Add(c)
					}

					ch <- pixel
					close(ch)
				}(j, i)
			}
		}
		close(results)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for ch := range results {
			pixel := <-ch
			writeColor(writer, pixel, ns)
		}
		wg.Done()
	}()

	wg.Wait()
}

func newCamera() Camera {
	var (
		lookfrom  = Point3{13, 2, 3}
		lookat    = Point3{0, 0, -0}
		vup       = Vec3{0, 1, 0}
		vfov      = 20.0
		aperture  = 0.1
		focusDist = 10.0
	)
	return NewCamera(
		imgWidth,
		imgHeight,
		lookfrom,
		lookat,
		vup,
		vfov,
		aperture,
		focusDist)
}

func main() {
	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	world := randomScene()
	cam := newCamera()

	// output image

	process(cam, world, imgWidth, imgHeight, samples, os.Stdout)
}
