package main

import "math"

// static values for Camera instances
const (
	// Aspect ratio
	ratio = 16.0 / 9.0

	// Focal length
	focalLen = 1.0
)

type Camera struct {
	origin, lowerLeftCorner Point3
	horiz, vert             Vec3
}

func NewCamera(lookfrom, lookat Point3, vup Vec3, vfov float64) Camera {
	var (
		// field of view
		theta      = vfov * (math.Pi / 180.0)
		h          = math.Tan(theta / 2)
		viewHeight = 2.0 * h
		viewWidth  = ratio * viewHeight

		// orientation
		w = lookfrom.Sub(lookat).Unit()
		u = vup.Cross(w).Unit()
		v = w.Cross(u)

		origin = lookfrom
		horiz  = u.MulS(viewWidth)
		vert   = v.MulS(viewHeight)
		llc    = origin.Sub(horiz.DivS(2), vert.DivS(2), w)
	)
	return Camera{origin, llc, horiz, vert}
}

func (c Camera) Ray(s, t float64) Ray {
	return Ray{c.origin, c.lowerLeftCorner.Add(c.horiz.MulS(s), c.vert.MulS(t)).Sub(c.origin)}
}
