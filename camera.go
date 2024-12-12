package main

import (
	"math"
)

type Camera struct {
	width, height           int
	lensRadius              float64
	origin, lowerLeftCorner Point3
	horiz, vert, u, v, w    Vec3
}

func NewCamera(width, height int, lookfrom, lookat Point3, vup Vec3, vfov, aperture, focusDist float64) Camera {
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
	return Camera{width, height, aperture / 2, origin, llc, horiz, vert, u, v, w}
}

func (c Camera) Ray(s, t float64) Ray {
	var (
		rd     = RandomVec3InUnitSphere().MulS(c.lensRadius)
		offset = c.u.MulS(rd.X).Add(c.v.MulS(rd.Y))
	)
	return Ray{c.origin.Add(offset), c.lowerLeftCorner.Add(c.horiz.MulS(s), c.vert.MulS(t)).Sub(c.origin, offset)}
}
