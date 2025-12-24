package main

import (
	"math"
	"testing"
)

func TestAABBHitRayFromOutsidePointingIn(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray from outside box, pointing at center
	ray := Ray{Orig: Point3{0, 0, -5}, Dir: Vec3{0, 0, 1}}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("ray from outside pointing at box should hit")
	}
}

func TestAABBHitRayFromOutsidePointingAway(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray from outside box, pointing away
	ray := Ray{Orig: Point3{0, 0, -5}, Dir: Vec3{0, 0, -1}}
	if box.Hit(ray, 0, 100) {
		t.Fatal("ray from outside pointing away should miss")
	}
}

func TestAABBHitRayFromInside(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray from inside box
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, 1}}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("ray from inside should hit box boundary")
	}
}

func TestAABBHitRayParallelToSide(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray parallel to X axis, inside the Y and Z slabs
	ray := Ray{Orig: Point3{-5, 0, 0}, Dir: Vec3{1, 0, 0}}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("ray parallel to X axis inside Y,Z slabs should hit")
	}

	// Ray parallel to X axis, outside the Y slab
	ray = Ray{Orig: Point3{-5, 2, 0}, Dir: Vec3{1, 0, 0}}
	if box.Hit(ray, 0, 100) {
		t.Fatal("ray parallel to X axis outside Y slab should miss")
	}
}

func TestAABBHitRayAtCorner(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray aimed at corner
	ray := Ray{Orig: Point3{-5, -5, -5}, Dir: Vec3{1, 1, 1}.Unit()}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("ray aimed at corner should hit")
	}
}

func TestAABBHitRayAtEdge(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray aimed at edge
	ray := Ray{Orig: Point3{-5, 0, 0}, Dir: Vec3{1, 0, 0}}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("ray aimed at edge should hit")
	}
}

func TestAABBHitTminTmaxRange(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})
	ray := Ray{Orig: Point3{0, 0, -5}, Dir: Vec3{0, 0, 1}}

	// Box is at t ≈ 4 to t ≈ 6
	// Should hit if range includes [4, 6]
	if !box.Hit(ray, 0, 100) {
		t.Fatal("should hit with wide range")
	}

	// Should miss if tmax is too small
	if box.Hit(ray, 0, 1) {
		t.Fatal("should miss if tmax < box distance")
	}

	// Should miss if tmin is too large
	if box.Hit(ray, 10, 100) {
		t.Fatal("should miss if tmin > box distance")
	}
}

func TestAABBHitGrazingRay(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray grazes the edge of the box (y = 1)
	ray := Ray{Orig: Point3{0, 1, -5}, Dir: Vec3{0, 0, 1}}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("grazing ray should hit")
	}
}

func TestAABBHitNegativeDirection(t *testing.T) {
	box := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	// Ray with negative direction components
	ray := Ray{Orig: Point3{5, 5, 5}, Dir: Vec3{-1, -1, -1}.Unit()}
	if !box.Hit(ray, 0, 100) {
		t.Fatal("ray with negative direction should hit")
	}
}

func TestAABBUnionEnclosesAllBoth(t *testing.T) {
	a := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})
	b := NewAABB(Point3{0.5, 0.5, 0.5}, Point3{2, 2, 2})

	union := a.Union(b)

	// Union should contain both boxes' min and max points
	if union.Min.X != -1 || union.Min.Y != -1 || union.Min.Z != -1 {
		t.Fatalf("union min = %v, want (-1, -1, -1)", union.Min)
	}
	if union.Max.X != 2 || union.Max.Y != 2 || union.Max.Z != 2 {
		t.Fatalf("union max = %v, want (2, 2, 2)", union.Max)
	}
}

func TestAABBUnionDisjoint(t *testing.T) {
	a := NewAABB(Point3{-2, -2, -2}, Point3{-1, -1, -1})
	b := NewAABB(Point3{1, 1, 1}, Point3{2, 2, 2})

	union := a.Union(b)

	// Union should span from a.Min to b.Max
	if union.Min.X != -2 || union.Min.Y != -2 || union.Min.Z != -2 {
		t.Fatalf("union min = %v, want (-2, -2, -2)", union.Min)
	}
	if union.Max.X != 2 || union.Max.Y != 2 || union.Max.Z != 2 {
		t.Fatalf("union max = %v, want (2, 2, 2)", union.Max)
	}
}

func TestAABBUnionNested(t *testing.T) {
	a := NewAABB(Point3{-2, -2, -2}, Point3{2, 2, 2})
	b := NewAABB(Point3{-1, -1, -1}, Point3{1, 1, 1})

	union := a.Union(b)

	// Union should equal the larger box (a)
	if union.Min != a.Min || union.Max != a.Max {
		t.Fatalf("union of nested boxes should equal larger box")
	}
}

func TestAABBHitLargeBox(t *testing.T) {
	// Very large box
	box := NewAABB(Point3{-1000, -1000, -1000}, Point3{1000, 1000, 1000})

	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{1, 0, 0}}
	if !box.Hit(ray, 0, math.MaxFloat64) {
		t.Fatal("ray from inside large box should hit")
	}
}
