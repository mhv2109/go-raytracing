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
}

func (h *HitRecord) setFaceNormal(r Ray, n Vec3) {
	h.F = r.Dir.Dot(n) < 0
	if h.F {
		h.N = n
	} else {
		h.N = n.Neg()
	}
}

type Hittable interface {
	// Hit checks if r interesects with Hittable. If so, it returns a HitRecord, nil otherwise.
	Hit(r Ray, tmin, tmax float64) *HitRecord
}

// Sphere is a shape defined by a Center point and a radius.
type Sphere struct {
	Center Point3
	R      float64
}

func (s Sphere) Hit(r Ray, tmin, tmax float64) *HitRecord {
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
		return nil
	}

	// Find the nearest root that lies in the acceptable range.
	var (
		sqrtd = math.Sqrt(d)
		root  = (-halfb - sqrtd) / a
	)

	if root < tmin || tmax < root {
		root = (-halfb + sqrtd) / a
		if root < tmin || tmax < root {
			return nil
		}
	}

	var (
		T = root
		P = r.At(T)
		N = P.Sub(s.Center).DivS(s.R)

		hr = HitRecord{P: P, T: T}
	)
	hr.setFaceNormal(r, N)
	return &hr
}

type Hittables struct {
	Objects []Hittable
}

func NewHittables(objects ...Hittable) Hittables {
	o := make([]Hittable, 0)
	h := Hittables{o}
	if len(objects) > 0 {
		h.Add(objects...)
	}
	return h
}

func (h *Hittables) Add(objects ...Hittable) {
	h.Objects = append(h.Objects, objects...)
}

func (h *Hittables) Clear() {
	h.Objects = make([]Hittable, 0)
}

func (h *Hittables) Hit(r Ray, tmin, tmax float64) *HitRecord {
	var (
		hr      *HitRecord
		closest = tmax
	)

	for _, object := range h.Objects {
		if temp := object.Hit(r, tmin, closest); temp != nil {
			closest = temp.T
			hr = temp
		}
	}

	return hr
}
