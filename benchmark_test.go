package main

import (
	"testing"
)

func newBenchmarkCamera() Camera {
	return NewCamera(
		defaultWidth/10,
		defaultHeight/10,
		defaultSamples/10,
		defaultDepth/10,
		Point3{13, 2, 3},
		Point3{0, 0, 0},
		Vec3{0, 1, 0},
		20.0,
		0.1,
		10.0)
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
