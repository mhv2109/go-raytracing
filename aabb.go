package main

import "math"

// AABB is an Axis-Aligned Bounding Box defined by minimum and maximum points.
// Used for accelerating ray-object intersection tests in the BVH.
type AABB struct {
	Min, Max Point3
}

// NewAABB creates an AABB from two corner points.
func NewAABB(min, max Point3) AABB {
	return AABB{Min: min, Max: max}
}

// Hit tests if a ray intersects the bounding box using the "slab method".
// This is the fastest known ray-box intersection algorithm.
//
// The slab method works by computing the intersection of the ray with each
// pair of parallel planes (x-min/x-max, y-min/y-max, z-min/z-max). The ray
// hits the box if and only if all three slab intervals overlap.
//
// Reference: "An Efficient and Robust Ray-Box Intersection Algorithm"
// (Williams et al., 2005)
func (a AABB) Hit(r Ray, tmin, tmax float64) bool {
	for axis := 0; axis < 3; axis++ {
		var (
			aMin, aMax, orig, dir float64
		)

		switch axis {
		case 0:
			aMin, aMax = a.Min.X, a.Max.X
			orig, dir = r.Orig.X, r.Dir.X
		case 1:
			aMin, aMax = a.Min.Y, a.Max.Y
			orig, dir = r.Orig.Y, r.Dir.Y
		case 2:
			aMin, aMax = a.Min.Z, a.Max.Z
			orig, dir = r.Orig.Z, r.Dir.Z
		}

		// Handle rays parallel to this axis (dir ≈ 0)
		if math.Abs(dir) < 1e-8 {
			// Ray is parallel to the slab planes
			// Check if ray origin is outside the slab
			if orig < aMin || orig > aMax {
				return false
			}
			// Ray is inside the slab for this axis, continue to next axis
			continue
		}

		// Compute intersection t values for this axis
		invD := 1.0 / dir
		t0 := (aMin - orig) * invD
		t1 := (aMax - orig) * invD

		// Swap if t0 > t1 (ray direction is negative on this axis)
		if t0 > t1 {
			t0, t1 = t1, t0
		}

		// Update the overall t range
		if t0 > tmin {
			tmin = t0
		}
		if t1 < tmax {
			tmax = t1
		}

		// Early exit if intervals don't overlap
		if tmin > tmax {
			return false
		}
	}

	return true
}

// Union returns an AABB that encompasses both this AABB and another.
func (a AABB) Union(b AABB) AABB {
	return AABB{
		Min: Point3{
			X: math.Min(a.Min.X, b.Min.X),
			Y: math.Min(a.Min.Y, b.Min.Y),
			Z: math.Min(a.Min.Z, b.Min.Z),
		},
		Max: Point3{
			X: math.Max(a.Max.X, b.Max.X),
			Y: math.Max(a.Max.Y, b.Max.Y),
			Z: math.Max(a.Max.Z, b.Max.Z),
		},
	}
}
