package main

import (
	"math"
	"math/rand"
	"slices"
)

// AABB is an axis-aligned bounding box defined by two corner points.
type AABB struct {
	Min, Max Point3
}

func (b AABB) Hit(r Ray, tmin, tmax float64) bool {
	// Slab method: check overlap of ray intervals on each axis.
	invD := 1.0 / r.Dir.X
	t0 := (b.Min.X - r.Orig.X) * invD
	t1 := (b.Max.X - r.Orig.X) * invD
	if invD < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}

	invD = 1.0 / r.Dir.Y
	t0 = (b.Min.Y - r.Orig.Y) * invD
	t1 = (b.Max.Y - r.Orig.Y) * invD
	if invD < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}

	invD = 1.0 / r.Dir.Z
	t0 = (b.Min.Z - r.Orig.Z) * invD
	t1 = (b.Max.Z - r.Orig.Z) * invD
	if invD < 0 {
		t0, t1 = t1, t0
	}
	if t0 > tmin {
		tmin = t0
	}
	if t1 < tmax {
		tmax = t1
	}
	if tmax <= tmin {
		return false
	}

	return true
}

// SurroundingBox returns the AABB that encloses both input boxes.
func SurroundingBox(a, b AABB) AABB {
	return AABB{
		Min: Vec3{
			math.Min(a.Min.X, b.Min.X),
			math.Min(a.Min.Y, b.Min.Y),
			math.Min(a.Min.Z, b.Min.Z),
		},
		Max: Vec3{
			math.Max(a.Max.X, b.Max.X),
			math.Max(a.Max.Y, b.Max.Y),
			math.Max(a.Max.Z, b.Max.Z),
		},
	}
}

// BVHNode is a node in a bounding volume hierarchy tree.
type BVHNode struct {
	Left, Right Hittable
	Box         AABB
}

// NewBVH builds a BVH tree from a slice of hittable objects.
func NewBVH(objects []Hittable) *BVHNode {
	n := &BVHNode{}

	axis := rand.Intn(3)
	cmp := func(a, b Hittable) int {
		ab := a.BoundingBox()
		bb := b.BoundingBox()
		var av, bv float64
		switch axis {
		case 0:
			av, bv = ab.Min.X, bb.Min.X
		case 1:
			av, bv = ab.Min.Y, bb.Min.Y
		default:
			av, bv = ab.Min.Z, bb.Min.Z
		}
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
		return 0
	}

	switch len(objects) {
	case 1:
		n.Left = objects[0]
		n.Right = objects[0]
	case 2:
		if cmp(objects[0], objects[1]) <= 0 {
			n.Left = objects[0]
			n.Right = objects[1]
		} else {
			n.Left = objects[1]
			n.Right = objects[0]
		}
	default:
		slices.SortFunc(objects, cmp)
		mid := len(objects) / 2
		n.Left = NewBVH(objects[:mid])
		n.Right = NewBVH(objects[mid:])
	}

	n.Box = SurroundingBox(n.Left.BoundingBox(), n.Right.BoundingBox())
	return n
}

func (n *BVHNode) Hit(r Ray, tmin, tmax float64, hr *HitRecord) bool {
	if !n.Box.Hit(r, tmin, tmax) {
		return false
	}

	hitLeft := n.Left.Hit(r, tmin, tmax, hr)
	if hitLeft {
		tmax = hr.T
	}
	hitRight := n.Right.Hit(r, tmin, tmax, hr)

	return hitLeft || hitRight
}

func (n *BVHNode) BoundingBox() AABB {
	return n.Box
}
