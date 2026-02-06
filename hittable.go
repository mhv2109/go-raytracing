package main

import "math"

// HitRecord captures the requisite details of a Ray intersecting with a Hittable.
type HitRecord struct {
	// Exact point of impact
	P Point3

	// Surface-normal vector
	N Vec3

	// Parameter t of impact
	T float64

	// "front" facing?
	F bool

	// Material of impacted object
	M Material
}

func NewHitRecord(P Point3, N Vec3, T float64, M Material, r Ray) HitRecord {
	hr := HitRecord{P: P, N: N, T: T, F: false, M: M}

	hr.F = r.Dir.Dot(N) < 0
	if !hr.F {
		hr.N = hr.N.Neg()
	}

	return hr
}

type Hittable interface {
	// Hit checks if r interesects with Hittable. If so, it returns a HitRecord, nil otherwise.
	Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool

	// BoundingBox returns the axis-aligned bounding box enclosing the hittable.
	BoundingBox() AABB
}

// Sphere is a shape defined by a Center point and a radius.
type Sphere struct {
	Center Point3
	R      float64
	M      Material
}

func (s Sphere) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	// A ray intersects the sphere if there exists two solutions for the quadratic
	// equation (P(t) - C) dot (P(t) - C) - r^2 = 0 for all t, where P(t) = A + t*halfb.
	// We can determine this by calulating the descriminant d. This has been
	// simplified using the method in section 6.2.
	var (
		oc = r.Orig.Sub(s.Center) // A - C

		a     = r.Dir.LenSq()
		halfb = oc.Dot(r.Dir)
		c     = oc.LenSq() - s.R*s.R

		d = halfb*halfb - a*c
	)

	if d < 0 {
		return false
	}

	// Find the nearest root that lies in the acceptable range.
	var (
		sqrtd = math.Sqrt(d)
		root  = (-halfb - sqrtd) / a
	)

	if root < tmin || tmax < root {
		root = (-halfb + sqrtd) / a
		if root < tmin || tmax < root {
			return false
		}
	}

	var (
		T    = root
		P    = r.At(T)
		N    = P.Sub(s.Center).DivS(s.R)
		temp = NewHitRecord(P, N, T, s.M, r)
	)
	*hr = temp
	return true
}

func (s Sphere) BoundingBox() AABB {
	offset := Vec3{s.R, s.R, s.R}
	return AABB{
		Min: s.Center.Sub(offset),
		Max: s.Center.Add(offset),
	}
}

type Hittables struct {
	Objects []Hittable
}

func NewHittables(objects ...Hittable) Hittables {
	o := make([]Hittable, len(objects))
	copy(o, objects)
	return Hittables{Objects: o}
}

func (h *Hittables) Add(objects ...Hittable) {
	h.Objects = append(h.Objects, objects...)
}

func (h *Hittables) Clear() {
	for i := range h.Objects {
		h.Objects[i] = nil
	}
	h.Objects = h.Objects[:0]
}

func (h *Hittables) BoundingBox() AABB {
	if len(h.Objects) == 0 {
		return AABB{}
	}
	box := h.Objects[0].BoundingBox()
	for _, obj := range h.Objects[1:] {
		box = SurroundingBox(box, obj.BoundingBox())
	}
	return box
}

func (h *Hittables) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	var (
		hit     = false
		closest = tmax
	)

	for _, object := range h.Objects {
		if object.Hit(r, tmin, closest, hr) {
			closest = hr.T
			hit = true
		}
	}

	return hit
}
