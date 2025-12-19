package main

import (
	"math"
	"testing"
)

func TestDiffusionScatterLambertian(t *testing.T) {
	albedo := Color{0.8, 0.3, 0.1}
	mat := NewDiffusion(albedo, WithDiffusionType(Lambertian))

	hr := HitRecord{
		P: Point3{0, 0, 0},
		N: Vec3{0, 0, 1},
		T: 1,
		F: true,
	}
	r := Ray{Orig: Point3{0, 0, -1}, Dir: Vec3{0, 0, 1}}

	var att Color
	var scatt Ray

	if !mat.Scatter(r, hr, &att, &scatt) {
		t.Fatalf("expected diffusion scatter to succeed")
	}

	if att != albedo {
		t.Fatalf("attenuation = %#v, want %#v", att, albedo)
	}

	if scatt.Dir.Len() == 0 {
		t.Fatalf("scattered direction has zero length")
	}

	// direction should be roughly in the hemisphere of the normal
	if scatt.Dir.Dot(hr.N) <= 0 {
		t.Fatalf("scattered direction not in hemisphere of normal: dot=%v", scatt.Dir.Dot(hr.N))
	}
}

func TestMetalScatterReflectsPerfectlyWithZeroFuzz(t *testing.T) {
	albedo := Color{0.8, 0.8, 0.8}
	metal := NewMetal(albedo, Fuzz(0))

	hr := HitRecord{
		P: Point3{0, 0, 0},
		N: Vec3{0, 0, 1},
		T: 1,
		F: true,
	}

	// Incoming ray at 45 degrees to the normal
	inDir := Vec3{0, 1, -1}.Unit()
	r := Ray{Orig: Point3{0, 0, 1}, Dir: inDir}

	var att Color
	var scatt Ray

	if !metal.Scatter(r, hr, &att, &scatt) {
		t.Fatalf("expected metal scatter to succeed")
	}

	if att != albedo {
		t.Fatalf("attenuation = %#v, want %#v", att, albedo)
	}

	wantDir := reflect(inDir, hr.N)
	if !vecAlmostEqual(scatt.Dir, wantDir) {
		t.Fatalf("reflected dir = %#v, want %#v", scatt.Dir, wantDir)
	}

	if scatt.Dir.Dot(hr.N) < 0 {
		t.Fatalf("reflected direction not above surface: dot=%v", scatt.Dir.Dot(hr.N))
	}
}

func TestDielectricScatterTotalInternalReflection(t *testing.T) {
	// Refractive index > 1, ray going from inside to outside at steep angle
	albedo := Color{1, 1, 1}
	d := NewDielectric(albedo, IndexOfRefraction(1.5))

	hr := HitRecord{
		P: Point3{0, 0, 0},
		N: Vec3{0, 0, 1},
		T: 1,
		F: false, // ray is inside, heading towards boundary
	}

	// Choose a direction such that ratio * sin(theta) > 1 for ratio=1.5
	angle := math.Pi / 3 // 60 degrees
	inDir := Vec3{math.Sin(angle), 0, math.Cos(angle)}.Unit()
	r := Ray{Orig: Point3{0, 0, -1}, Dir: inDir}

	var att Color
	var scatt Ray

	if !d.Scatter(r, hr, &att, &scatt) {
		t.Fatalf("expected dielectric scatter to succeed")
	}

	// Total internal reflection should produce a reflected ray
	reflDir := reflect(inDir, hr.N)

	// Check that scattered direction is close to reflection direction
	if scatt.Dir.Dot(reflDir) < 0.999 {
		t.Fatalf("expected total internal reflection, dot(scatt, refl)=%v", scatt.Dir.Dot(reflDir))
	}
}

func TestDielectricScatterRefraction(t *testing.T) {
	albedo := Color{1, 1, 1}
	d := NewDielectric(albedo, IndexOfRefraction(1.5))

	hr := HitRecord{
		P: Point3{0, 0, 0},
		N: Vec3{0, 0, 1},
		T: 1,
		F: true, // ray coming from air into dielectric
	}

	inDir := Vec3{0, 0, -1} // straight on
	r := Ray{Orig: Point3{0, 0, 1}, Dir: inDir}

	var att Color
	var scatt Ray

	if !d.Scatter(r, hr, &att, &scatt) {
		t.Fatalf("expected dielectric scatter to succeed")
	}

	// For normal incidence, refraction direction should align with -normal (into the material)
	if scatt.Dir.Dot(hr.N.Neg()) <= 0.999 {
		t.Fatalf("expected refracted direction into material, dot=%v", scatt.Dir.Dot(hr.N.Neg()))
	}
}
