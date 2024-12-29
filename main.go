package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/schollz/progressbar/v3"
)

var (
	// cmdline args
	imgWidth   int
	imgHeight  int
	samples    int
	depth      int
	jobs       int
	simpleDiff bool
	cpuprofile string
	outputFile string

	// defaults
	defaultWidth   = 2560
	defaultHeight  = 1440
	defaultSamples = 500
	defaultDepth   = 50
	defaultJobs    = 16 * runtime.NumCPU() // determined by benchmarking
)

func init() {
	flag.IntVar(&imgWidth, "width", 2560, "output image width")
	flag.IntVar(&imgHeight, "height", 1440, "output image height")
	flag.IntVar(&samples, "samples", 500, "number of samples per pixel")
	flag.IntVar(&depth, "depth", 50, "number of ray bounces to calculate")
	flag.IntVar(&jobs, "jobs", defaultJobs, "number of jobs for rendering")
	flag.BoolVar(&simpleDiff, "simple", false, "use simple diffusion calculation")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "create a CPU profile and save to file")
	flag.StringVar(&outputFile, "output", "", "output file, defaults to stdout")
}

// diffustionMaterial allows us to select the diffusion function at runtime
func diffusionMaterial() DiffusionOpt {
	if simpleDiff {
		return WithDiffusionType(SimpleDiffusion)
	}
	return WithDiffusionType(Lambertian)
}

func randomScene() Hittables {
	world := NewHittables()

	// earth/ground/floor

	ground := Sphere{
		Point3{0, -1000, 0},
		1000,
		NewDiffusion(Color{0.5, 0.5, 0.5}, diffusionMaterial()),
	}
	world.Add(ground)

	// big spheres in the center

	sphere1 := Sphere{
		Point3{0, 1, 0},
		1,
		NewDielectric(Color{1, 1, 1}, IndexOfRefraction(1.5)),
	}
	world.Add(sphere1)

	sphere2 := Sphere{
		Point3{-4, 1, 0},
		1,
		NewDiffusion(Color{0.4, 0.2, 0.1}, diffusionMaterial()),
	}
	world.Add(sphere2)

	sphere3 := Sphere{
		Point3{4, 1, 0},
		1,
		NewMetal(Color{0.7, 0.6, 0.5}, Fuzz(0)),
	}
	world.Add(sphere3)

	// add random little spheres all over the ground

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Point3{float64(a) + 0.8*rand.Float64(), 0.2, float64(b) + 0.8*rand.Float64()}
			if (center.Sub(sphere1.Center).Len() > 1.2) &&
				(center.Sub(sphere2.Center).Len() > 1.2) &&
				(center.Sub(sphere3.Center).Len() > 1.2) {

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
					c = Color{1, 1, 1}
					m = NewDielectric(c, IndexOfRefraction(1.5))
				}

				world.Add(Sphere{center, 0.2, m})
			}
		}
	}

	return world
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
		samples,
		depth,
		jobs,
		lookfrom,
		lookat,
		vup,
		vfov,
		aperture,
		focusDist)
}

func main() {
	flag.Parse()

	// configure profiling

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

	// setup output

	output := os.Stdout
	if outputFile != "" {
		var err error
		output, err = os.Create(outputFile)
		if err != nil {
			log.Fatal("could not create output file: ", err)
		}
		defer output.Close()
	}

	// output image

	cam := newCamera()

	fmt.Fprintln(output, "P3")
	fmt.Fprintln(output, cam.ImageWidth(), cam.ImageHeight())
	fmt.Fprintln(output, "255")

	bar := progressbar.Default(int64(cam.ImageSize()))
	for pixel := range cam.Render(randomScene()) {
		fmt.Fprintln(output, pixel.R, pixel.G, pixel.B)
		bar.Add(1)
	}
}
