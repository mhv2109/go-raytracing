package main

// static values for Camera instances
const (
	// Aspect ratio
	ratio = 16.0 / 9.0

	// Focal length and field of view
	// Since these are "static" values, you can think of this as a "prime" lens attached to each Camera instance
	viewHeight = 2.0
	viewWidth  = ratio * viewHeight
	focalLen   = 1.0
)

type Camera struct {
	origin, lowerLeftCorner Point3
	horiz, vert             Vec3
}

func NewCamera(origin Point3) Camera {
	var (
		horiz = Vec3{viewWidth, 0, 0}
		vert  = Vec3{0, viewHeight, 0}
		llc   = origin.Sub(horiz.DivS(2), vert.DivS(2), Point3{0, 0, focalLen})
	)
	return Camera{origin, llc, horiz, vert}
}

func (c Camera) Ray(u, v float64) Ray {
	return Ray{c.origin, c.lowerLeftCorner.Add(c.horiz.MulS(u), c.vert.MulS(v), c.origin.Neg())}
}
