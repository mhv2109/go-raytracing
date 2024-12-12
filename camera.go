package main

import (
	"fmt"
	"io"
	"math"
	"sync"
)

type Camera struct {
	width, height           int
	samples, depth          int
	jobs                    int
	lensRadius              float64
	origin, lowerLeftCorner Point3
	horiz, vert, u, v, w    Vec3
}

func NewCamera(width, height, samples, depth, jobs int, lookfrom, lookat Point3, vup Vec3, vfov, aperture, focusDist float64) Camera {
	var (
		// field of view
		theta      = vfov * (math.Pi / 180.0)
		h          = math.Tan(theta / 2)
		viewHeight = 2.0 * h
		ratio      = float64(width) / float64(height)
		viewWidth  = ratio * viewHeight

		// orientation
		w = lookfrom.Sub(lookat).Unit()
		u = vup.Cross(w).Unit()
		v = w.Cross(u)

		origin = lookfrom
		horiz  = u.MulS(viewWidth).MulS(focusDist)
		vert   = v.MulS(viewHeight).MulS(focusDist)
		llc    = origin.Sub(horiz.DivS(2), vert.DivS(2), w.MulS(focusDist))
	)
	return Camera{width, height, samples, depth, jobs, aperture / 2, origin, llc, horiz, vert, u, v, w}
}

func (cam Camera) ray(s, t float64) Ray {
	var (
		rd     = RandomVec3InUnitSphere().MulS(cam.lensRadius)
		offset = cam.u.MulS(rd.X).Add(cam.v.MulS(rd.Y))
	)
	return Ray{
		cam.origin.Add(offset),
		cam.lowerLeftCorner.Add(cam.horiz.MulS(s), cam.vert.MulS(t)).Sub(cam.origin, offset),
	}
}

func (Camera) clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

func (cam Camera) writeColor(w io.Writer, c Color, samples int) {
	// Divide the color by the number of samples and scale float values [0, 1]
	// to [0, 255]
	var (
		scale = 1.0 / float64(samples)
		r     = int(255.999 * cam.clamp(math.Sqrt(c.X*scale), 0.0, 0.999))
		g     = int(255.999 * cam.clamp(math.Sqrt(c.Y*scale), 0.0, 0.999))
		b     = int(255.999 * cam.clamp(math.Sqrt(c.Z*scale), 0.0, 0.999))
	)
	fmt.Fprintln(w, r, g, b)
}

// rayColor calculates the Color along the Ray. We define objects + colors here,
// and return an object's color if the Ray intersects it. Otherwise, we return
// the background color
func (cam Camera) rayColor(r Ray, world Hittables) Color {
	var (
		mult  = Vec3{1, 1, 1}
		hr    HitRecord
		n     = 0
		att   Color
		scatt Ray
	)

LOOP:
	if n > cam.depth {
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

func (cam Camera) Render(world Hittables, writer io.Writer) {
	fmt.Fprintln(writer, "P3")
	fmt.Fprintln(writer, cam.width, cam.height)
	fmt.Fprintln(writer, "255")

	// Pan across each pixel of the output image and calculate the color of each.
	var (
		wg      sync.WaitGroup
		results = make(chan chan Color, cam.jobs)
		rands   = make([]float64, 2*cam.samples)
	)

	// calculate random values up-front
	for i := range rands {
		rands[i] = <-RandomCh
	}

	wg.Add(1)
	go func() {
		for j := cam.height; j >= 0; j-- {
			for i := 0; i < cam.width; i++ {
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

					for s := 0; s < 2*cam.samples; s += 2 {
						u = (float64(i) + rands[s]) / (float64(cam.width) - 1)
						v = (float64(j) + rands[s+1]) / (float64(cam.height) - 1)
						r = cam.ray(u, v)
						c = cam.rayColor(r, world)
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
			cam.writeColor(writer, pixel, cam.samples)
		}
		wg.Done()
	}()

	wg.Wait()
}
