package main

import (
	"math"
	"math/rand"
	"testing"
)

func newBenchmarkCamera() Camera {
	return NewCamera(
		defaultWidth/10,
		defaultHeight/10,
		defaultSamples/10,
		defaultDepth/10,
		defaultJobs,
		Point3{13, 2, 3},
		Point3{0, 0, 0},
		Vec3{0, 1, 0},
		20.0,
		0.1,
		10.0)
}

// createRandomScene creates the scene objects without building BVH or Hittables
func createRandomSceneObjects() []Hittable {
	objects := make([]Hittable, 0, 500)

	// Ground
	objects = append(objects, Sphere{
		Point3{0, -1000, 0},
		1000,
		NewDiffusion(Color{0.5, 0.5, 0.5}, diffusionMaterial()),
	})

	// Big spheres
	sphere1 := Sphere{Point3{0, 1, 0}, 1, NewDielectric(Color{1, 1, 1}, IndexOfRefraction(1.5))}
	sphere2 := Sphere{Point3{-4, 1, 0}, 1, NewDiffusion(Color{0.4, 0.2, 0.1}, diffusionMaterial())}
	sphere3 := Sphere{Point3{4, 1, 0}, 1, NewMetal(Color{0.7, 0.6, 0.5}, Fuzz(0))}

	objects = append(objects, sphere1, sphere2, sphere3)

	// Random small spheres
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Point3{float64(a) + 0.8*rand.Float64(), 0.2, float64(b) + 0.8*rand.Float64()}
			if (center.Sub(sphere1.Center).Len() > 1.2) &&
				(center.Sub(sphere2.Center).Len() > 1.2) &&
				(center.Sub(sphere3.Center).Len() > 1.2) {

				var m Material
				switch choose := rand.Float64(); {
				case choose < 0.8:
					c := RandomVec3(0, 1).Mul(RandomVec3(0, 1))
					m = NewDiffusion(c, diffusionMaterial())
				case choose < 0.95:
					c := RandomVec3(0.5, 1)
					fuzz := rand.Float64() * 0.5
					m = NewMetal(c, Fuzz(fuzz))
				default:
					m = NewDielectric(Color{1, 1, 1}, IndexOfRefraction(1.5))
				}

				objects = append(objects, Sphere{center, 0.2, m})
			}
		}
	}

	return objects
}

func BenchmarkRender(b *testing.B) {
	var (
		world = randomScene()
		cam   = newBenchmarkCamera()
	)

	b.Run("render", func(b *testing.B) {
		for _ = range cam.Render(world) {
		}
	})
}

func BenchmarkBVHConstruction(b *testing.B) {
	objects := createRandomSceneObjects()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBVH(objects)
	}
}

func BenchmarkHitLinear(b *testing.B) {
	objects := createRandomSceneObjects()
	world := NewHittables(objects...)

	ray := Ray{Orig: Point3{13, 2, 3}, Dir: Vec3{-1, -0.1, -0.2}.Unit()}
	var hr HitRecord

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Hit(ray, 0.001, math.MaxFloat64, &hr)
	}
}

func BenchmarkHitBVH(b *testing.B) {
	objects := createRandomSceneObjects()
	bvh := NewBVH(objects)

	ray := Ray{Orig: Point3{13, 2, 3}, Dir: Vec3{-1, -0.1, -0.2}.Unit()}
	var hr HitRecord

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bvh.Hit(ray, 0.001, math.MaxFloat64, &hr)
	}
}

func BenchmarkRenderLinear(b *testing.B) {
	objects := createRandomSceneObjects()
	world := NewHittables(objects...)
	cam := newBenchmarkCamera()

	b.ResetTimer()
	for _ = range cam.Render(&world) {
	}
}

func BenchmarkRenderBVH(b *testing.B) {
	objects := createRandomSceneObjects()
	bvh := NewBVH(objects)
	cam := newBenchmarkCamera()

	b.ResetTimer()
	for _ = range cam.Render(bvh) {
	}
}
