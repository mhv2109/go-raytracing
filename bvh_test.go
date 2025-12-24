package main

import (
	"math"
	"testing"
)

func TestBVHSingleObject(t *testing.T) {
	sphere := Sphere{Center: Point3{0, 0, -1}, R: 0.5, M: NewDiffusion(Color{1, 1, 1})}
	bvh := NewBVH([]Hittable{sphere})

	// Ray hitting the sphere
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, -1}}
	var hr HitRecord

	if !bvh.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("BVH should hit sphere")
	}

	// Verify hit distance matches sphere hit
	if !almostEqual(hr.T, 0.5) {
		t.Fatalf("hr.T = %v, want ~0.5", hr.T)
	}
}

func TestBVHTwoObjects(t *testing.T) {
	near := Sphere{Center: Point3{0, 0, -1}, R: 0.5, M: NewDiffusion(Color{1, 1, 1})}
	far := Sphere{Center: Point3{0, 0, -3}, R: 0.5, M: NewDiffusion(Color{1, 1, 1})}
	bvh := NewBVH([]Hittable{near, far})

	// Ray going through both spheres
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, -1}}
	var hr HitRecord

	if !bvh.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("BVH should hit")
	}

	// Should hit nearest sphere first
	if !almostEqual(hr.T, 0.5) {
		t.Fatalf("nearest hit T = %v, want ~0.5", hr.T)
	}
}

func TestBVHChoosesNearest(t *testing.T) {
	// Create spheres at different distances
	near := Sphere{Center: Point3{0, 0, -2}, R: 0.5, M: NewDiffusion(Color{1, 0, 0})}
	mid := Sphere{Center: Point3{0, 0, -5}, R: 0.5, M: NewDiffusion(Color{0, 1, 0})}
	far := Sphere{Center: Point3{0, 0, -8}, R: 0.5, M: NewDiffusion(Color{0, 0, 1})}

	// Build BVH with objects in random order
	bvh := NewBVH([]Hittable{far, near, mid})

	// Ray going through all three
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, -1}}
	var hr HitRecord

	if !bvh.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("BVH should hit")
	}

	// Should hit nearest sphere (at z=-2, radius=0.5, so t=1.5)
	if !almostEqual(hr.T, 1.5) {
		t.Fatalf("nearest hit T = %v, want ~1.5", hr.T)
	}
}

func TestBVHMiss(t *testing.T) {
	sphere := Sphere{Center: Point3{0, 0, -1}, R: 0.5, M: NewDiffusion(Color{1, 1, 1})}
	bvh := NewBVH([]Hittable{sphere})

	// Ray missing the sphere
	ray := Ray{Orig: Point3{2, 0, 0}, Dir: Vec3{0, 0, -1}}
	var hr HitRecord

	if bvh.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("BVH should miss - ray doesn't hit sphere")
	}
}

func TestBVHManyObjectsMatchesLinear(t *testing.T) {
	// Create a grid of spheres
	objects := make([]Hittable, 0, 100)
	for x := -5; x < 5; x++ {
		for z := -5; z < 5; z++ {
			sphere := Sphere{
				Center: Point3{float64(x), 0, float64(z) - 10},
				R:      0.3,
				M:      NewDiffusion(Color{1, 1, 1}),
			}
			objects = append(objects, sphere)
		}
	}

	bvh := NewBVH(objects)
	linear := NewHittables(objects...)

	// Test with many random rays
	for i := 0; i < 100; i++ {
		ray := Ray{
			Orig: Point3{0, 0, 0},
			Dir:  RandomVec3(-1, 1).Unit(),
		}

		var hr1, hr2 HitRecord
		hit1 := bvh.Hit(ray, 0.001, math.MaxFloat64, &hr1)
		hit2 := linear.Hit(ray, 0.001, math.MaxFloat64, &hr2)

		if hit1 != hit2 {
			t.Fatalf("ray %d: hit mismatch: BVH=%v, linear=%v", i, hit1, hit2)
		}

		if hit1 && !almostEqual(hr1.T, hr2.T) {
			t.Fatalf("ray %d: T mismatch: BVH=%v, linear=%v", i, hr1.T, hr2.T)
		}

		if hit1 && !vecAlmostEqual(hr1.P, hr2.P) {
			t.Fatalf("ray %d: P mismatch: BVH=%v, linear=%v", i, hr1.P, hr2.P)
		}
	}
}

func TestBVHBoundingBox(t *testing.T) {
	s1 := Sphere{Center: Point3{0, 0, 0}, R: 1, M: NewDiffusion(Color{1, 1, 1})}
	s2 := Sphere{Center: Point3{3, 0, 0}, R: 1, M: NewDiffusion(Color{1, 1, 1})}

	bvh := NewBVH([]Hittable{s1, s2})
	box := bvh.BoundingBox()

	// BVH bounding box should contain both spheres
	// s1 spans [-1, 1] in all dimensions
	// s2 center at (3, 0, 0) with R=1 spans [2, 4] in X, [-1, 1] in Y, Z
	// Union should be [-1, 4] in X, [-1, 1] in Y, [-1, 1] in Z

	if !almostEqual(box.Min.X, -1) || !almostEqual(box.Max.X, 4) {
		t.Fatalf("box X range = [%v, %v], want [-1, 4]", box.Min.X, box.Max.X)
	}
	if !almostEqual(box.Min.Y, -1) || !almostEqual(box.Max.Y, 1) {
		t.Fatalf("box Y range = [%v, %v], want [-1, 1]", box.Min.Y, box.Max.Y)
	}
	if !almostEqual(box.Min.Z, -1) || !almostEqual(box.Max.Z, 1) {
		t.Fatalf("box Z range = [%v, %v], want [-1, 1]", box.Min.Z, box.Max.Z)
	}
}

func TestBVHEarlyRejection(t *testing.T) {
	// Create spheres far away
	objects := make([]Hittable, 0, 10)
	for i := 0; i < 10; i++ {
		sphere := Sphere{
			Center: Point3{float64(i) * 10, 0, -100},
			R:      1,
			M:      NewDiffusion(Color{1, 1, 1}),
		}
		objects = append(objects, sphere)
	}

	bvh := NewBVH(objects)

	// Ray pointing completely away from spheres
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 1, 0}}
	var hr HitRecord

	if bvh.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("BVH should reject ray early via bounding box")
	}
}

func TestBVHLargeScene(t *testing.T) {
	// Create a large scene similar to randomScene
	objects := make([]Hittable, 0, 500)

	// Ground sphere
	objects = append(objects, Sphere{
		Center: Point3{0, -1000, 0},
		R:      1000,
		M:      NewDiffusion(Color{0.5, 0.5, 0.5}),
	})

	// Add many small spheres (but leave some gaps)
	for a := -10; a < 10; a++ {
		for b := -10; b < 10; b++ {
			// Skip center to leave a gap
			if a == 0 && b == 0 {
				continue
			}
			center := Point3{float64(a), 0.2, float64(b)}
			objects = append(objects, Sphere{
				Center: center,
				R:      0.2,
				M:      NewDiffusion(Color{0.5, 0.5, 0.5}),
			})
		}
	}

	// Build BVH
	bvh := NewBVH(objects)

	// Test a ray through the gap that should hit the ground
	ray := Ray{Orig: Point3{0, 1, 0}, Dir: Vec3{0, -1, 0}}
	var hr HitRecord

	if !bvh.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("should hit ground sphere")
	}

	// Verify we hit ground (large sphere at y=-1000, radius=1000)
	// Ray from (0,1,0) going down hits ground at y=0 (top of ground sphere)
	if !almostEqual(hr.P.Y, 0) {
		t.Fatalf("expected to hit ground at y=0, got P.Y=%v", hr.P.Y)
	}
}
