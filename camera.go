package main

import (
	"iter"
	"math"
	"math/rand"
	"sync"
)

type Camera struct {
	width, height           int
	samples, depth          int
	lensRadius              float64
	origin, lowerLeftCorner Point3
	horiz, vert, u, v, w    Vec3
}

func NewCamera(width, height, samples, depth int, lookfrom, lookat Point3, vup Vec3, vfov, aperture, focusDist float64) Camera {
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
	return Camera{width, height, samples, depth, aperture / 2, origin, llc, horiz, vert, u, v, w}
}

func (cam Camera) ImageWidth() int {
	return cam.width
}

func (cam Camera) ImageHeight() int {
	return cam.height
}

func (cam Camera) ImageSize() int {
	return cam.ImageWidth() * cam.ImageHeight()
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

func (cam Camera) coords() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for j := cam.height; j >= 0; j-- {
			for i := 0; i < cam.width; i++ {
				if !yield(j, i) {
					return
				}
			}
		}
	}
}

func (cam Camera) renderPixel(world Hittables, j, i int) RGB {
	var (
		u, v    float64
		r       Ray
		c       Color
		wg      sync.WaitGroup
		samples = make([]Color, cam.samples)
	)

	for s := 0; s < cam.samples; s++ {
		wg.Add(1)
		go func(s int) {
			u = (float64(i) + rand.Float64()) / (float64(cam.width) - 1)
			v = (float64(j) + rand.Float64()) / (float64(cam.height) - 1)
			r = cam.ray(u, v)
			c = cam.rayColor(r, world)
			samples[s] = c
			wg.Done()
		}(s)
	}

	wg.Wait()

	return samples[0].Add(samples[1:]...).RGB(float64(cam.samples))
}

func (cam Camera) render(world Hittables, coords iter.Seq2[int, int]) iter.Seq[RGB] {
	return func(yield func(RGB) bool) {
		for j, i := range coords {
			if !yield(cam.renderPixel(world, j, i)) {
				return
			}
		}
	}
}

func (cam Camera) Render(world Hittables) iter.Seq[RGB] {
	// Pan across each pixel of the output image and calculate the color of each.
	return cam.render(world, cam.coords())
}
