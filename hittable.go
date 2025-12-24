package main

import (
	"math"
	"math/rand"
	"sort"
)

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
	// BoundingBox returns an axis-aligned bounding box for the hittable between time0 and time1.
	// For this static scene, time parameters are ignored but retained for API flexibility.
	BoundingBox(time0, time1 float64, box *AABB) bool
}

// Sphere is a shape defined by a Center point and a radius.
type Sphere struct {
	Center Point3
	R      float64
	M      Material
}

func (s Sphere) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	var (
		oc = r.Orig.Sub(s.Center)

		a     = r.Dir.LenSq()
		halfb = oc.Dot(r.Dir)
		c     = oc.LenSq() - s.R*s.R

		d = halfb*halfb - a*c
	)

	if d < 0 {
		return false
	}

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

func (s Sphere) BoundingBox(time0, time1 float64, box *AABB) bool {
	radius := Vec3{s.R, s.R, s.R}
	*box = AABB{
		Min: s.Center.Sub(radius),
		Max: s.Center.Add(radius),
	}
	return true
}

type Hittables struct {
	Objects []Hittable
	bvhRoot Hittable
}

func NewHittables(objects ...Hittable) Hittables {
	o := make([]Hittable, len(objects))
	copy(o, objects)
	return Hittables{Objects: o}
}

func (h *Hittables) Add(objects ...Hittable) {
	h.Objects = append(h.Objects, objects...)
	h.bvhRoot = nil
}

func (h *Hittables) Clear() {
	for i := range h.Objects {
		h.Objects[i] = nil
	}
	h.Objects = h.Objects[:0]
	h.bvhRoot = nil
}

type BVHNode struct {
	left, right Hittable
	box        AABB
}

func NewBVHNode(objects []Hittable, time0, time1 float64) Hittable {
	n := len(objects)
	if n == 0 {
		return nil
	}
	if n == 1 {
		var box AABB
		if !objects[0].BoundingBox(time0, time1, &box) {
			return nil
		}
		return &BVHNode{left: objects[0], right: nil, box: box}
	}

	axis := rand.Intn(3)
	less := func(i, j int) bool {
		var boxI, boxJ AABB
		if !objects[i].BoundingBox(time0, time1, &boxI) || !objects[j].BoundingBox(time0, time1, &boxJ) {
			return false
		}
		switch axis {
		case 0:
			return boxI.Min.X < boxJ.Min.X
		case 1:
			return boxI.Min.Y < boxJ.Min.Y
		default:
			return boxI.Min.Z < boxJ.Min.Z
		}
	}

	sort.Slice(objects, less)

	mid := n / 2
	left := NewBVHNode(objects[:mid], time0, time1)
	right := NewBVHNode(objects[mid:], time0, time1)

	var boxLeft, boxRight AABB
	if left == nil || right == nil || !left.BoundingBox(time0, time1, &boxLeft) || !right.BoundingBox(time0, time1, &boxRight) {
		return nil
	}

	return &BVHNode{
		left:  left,
		right: right,
		box:  SurroundingBox(boxLeft, boxRight),
	}
}

func (n *BVHNode) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	if n == nil {
		return false
	}
	if !n.box.Hit(r, tmin, tmax) {
		return false
	}

	hitLeft := n.left != nil && n.left.Hit(r, tmin, tmax, hr)
	if hitLeft {
		tmax = hr.T
	}
	hitRight := n.right != nil && n.right.Hit(r, tmin, tmax, hr)

	return hitLeft || hitRight
}

func (n *BVHNode) BoundingBox(time0, time1 float64, box *AABB) bool {
	if n == nil {
		return false
	}
	*box = n.box
	return true
}

func SurroundingBox(box0, box1 AABB) AABB {
	small := Point3{
		X: math.Min(box0.Min.X, box1.Min.X),
		Y: math.Min(box0.Min.Y, box1.Min.Y),
		Z: math.Min(box0.Min.Z, box1.Min.Z),
	}
	big := Point3{
		X: math.Max(box0.Max.X, box1.Max.X),
		Y: math.Max(box0.Max.Y, box1.Max.Y),
		Z: math.Max(box0.Max.Z, box1.Max.Z),
	}
	return AABB{Min: small, Max: big}
}

func (h *Hittables) BuildBVH() {
	if len(h.Objects) == 0 {
		h.bvhRoot = nil
		return
	}
	objs := make([]Hittable, len(h.Objects))
	copy(objs, h.Objects)
	h.bvhRoot = NewBVHNode(objs, 0, 0)
}

func (h *Hittables) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	if h.bvhRoot != nil {
		return h.bvhRoot.Hit(r, tmin, tmax, hr)
	}

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
