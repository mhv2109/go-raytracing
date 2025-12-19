package main

import (
	"math"
	"testing"
)

const floatEps = 1e-6

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < floatEps
}

func vecAlmostEqual(a, b Vec3) bool {
	return almostEqual(a.X, b.X) && almostEqual(a.Y, b.Y) && almostEqual(a.Z, b.Z)
}

func TestHittablesHitEmpty(t *testing.T) {
	var world Hittables

	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{1, 0, 0}}
	var hr HitRecord

	if hit := world.Hit(ray, 0.001, math.MaxFloat64, &hr); hit {
		t.Fatalf("expected no hit on empty Hittables, got hit = %v", hit)
	}
}

func TestSphereSingleHit(t *testing.T) {
	sphere := Sphere{Center: Point3{0, 0, -1}, R: 0.5}
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, -1}}

	var hr HitRecord
	if !sphere.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatalf("expected ray to hit sphere")
	}

	// Analytical intersection for sphere at (0,0,-1) radius 0.5 and ray from origin along -Z
	// t should be close to 0.5
	if !almostEqual(hr.T, 0.5) {
		t.Fatalf("hr.T = %v, want ~0.5", hr.T)
	}

	wantP := Point3{0, 0, -0.5}
	if !vecAlmostEqual(Vec3(hr.P), Vec3(wantP)) {
		t.Fatalf("hit point = %#v, want %#v", hr.P, wantP)
	}

	wantN := Vec3{0, 0, 1}
	if !vecAlmostEqual(hr.N, wantN) {
		t.Fatalf("normal = %#v, want %#v", hr.N, wantN)
	}
}

func TestHittablesHitChoosesNearest(t *testing.T) {
	near := Sphere{Center: Point3{0, 0, -1}, R: 0.5}
	far := Sphere{Center: Point3{0, 0, -3}, R: 0.5}

	world := NewHittables(near, far)
	var hr HitRecord
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, -1}}

	if !world.Hit(ray, 0.001, math.MaxFloat64, &hr) {
		t.Fatalf("expected ray to hit world")
	}

	// Nearest intersection is with the "near" sphere at t ~= 0.5
	if !almostEqual(hr.T, 0.5) {
		t.Fatalf("nearest hit T = %v, want ~0.5", hr.T)
	}

	wantP := Point3{0, 0, -0.5}
	if !vecAlmostEqual(Vec3(hr.P), Vec3(wantP)) {
		t.Fatalf("hit point = %#v, want %#v", hr.P, wantP)
	}
}

func TestNewHitRecordFrontFace(t *testing.T) {
	P := Point3{0, 0, -1}
	N := Vec3{0, 0, 1}
	ray := Ray{Orig: Point3{0, 0, 0}, Dir: Vec3{0, 0, -1}}

	hr := NewHitRecord(P, N, 1.0, nil, ray)

	if !hr.F {
		t.Fatalf("expected front face, got F = false")
	}

	if !vecAlmostEqual(hr.N, N) {
		t.Fatalf("normal = %#v, want %#v", hr.N, N)
	}
}

func TestNewHitRecordBackFace(t *testing.T) {
	P := Point3{0, 0, -1}
	N := Vec3{0, 0, 1}
	// Ray going in same direction as normal means we are inside the surface
	ray := Ray{Orig: Point3{0, 0, -1}, Dir: Vec3{0, 0, 1}}

	hr := NewHitRecord(P, N, 1.0, nil, ray)

	if hr.F {
		t.Fatalf("expected back face, got F = true")
	}

	wantN := N.Neg()
	if !vecAlmostEqual(hr.N, wantN) {
		t.Fatalf("normal = %#v, want %#v", hr.N, wantN)
	}
}
