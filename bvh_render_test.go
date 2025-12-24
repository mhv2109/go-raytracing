package main

import (
	"math"
	"testing"
)

func TestBVHWithRandomScene(t *testing.T) {
	// Create actual random scene
	bvh := randomScene()
	linear := NewHittables(createRandomSceneObjects()...)

	// Test rays from the camera position
	testRays := []Ray{
		// Ray toward center
		{Orig: Point3{13, 2, 3}, Dir: Vec3{-1, -0.1, -0.2}.Unit()},
		// Ray toward ground
		{Orig: Point3{0, 5, 0}, Dir: Vec3{0, -1, 0}},
		// Ray at angle
		{Orig: Point3{13, 2, 3}, Dir: Vec3{-0.5, 0, -0.5}.Unit()},
	}

	for i, ray := range testRays {
		var hr1, hr2 HitRecord
		hit1 := bvh.Hit(ray, 0.001, math.MaxFloat64, &hr1)
		hit2 := linear.Hit(ray, 0.001, math.MaxFloat64, &hr2)

		if hit1 != hit2 {
			t.Fatalf("Ray %d: hit mismatch: BVH=%v, linear=%v", i, hit1, hit2)
		}

		if hit1 {
			if !almostEqual(hr1.T, hr2.T) {
				t.Fatalf("Ray %d: T mismatch: BVH=%v, linear=%v", i, hr1.T, hr2.T)
			}
			t.Logf("Ray %d: Hit at T=%v", i, hr1.T)
		} else {
			t.Logf("Ray %d: Miss", i)
		}
	}
}

func TestAABBWithActualSceneBounds(t *testing.T) {
	// Test AABB with the actual ground sphere
	ground := Sphere{
		Center: Point3{0, -1000, 0},
		R:      1000,
		M:      NewDiffusion(Color{0.5, 0.5, 0.5}),
	}

	bbox := ground.BoundingBox()
	t.Logf("Ground bbox: min=%v, max=%v", bbox.Min, bbox.Max)

	// Ray from above going down
	ray := Ray{Orig: Point3{0, 5, 0}, Dir: Vec3{0, -1, 0}}

	// Test AABB hit
	if !bbox.Hit(ray, 0.001, math.MaxFloat64) {
		t.Fatal("AABB should hit ground sphere bbox")
	}

	// Test actual sphere hit
	var hr HitRecord
	if !ground.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatal("Ray should hit ground sphere")
	}
	t.Logf("Ground hit at T=%v, P=%v", hr.T, hr.P)
}
