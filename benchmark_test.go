package main

import (
	"runtime"
	"testing"
)

type writer struct{}

func (w writer) Write(p []byte) (int, error) {
	return len(p), nil
}

func newBenchmarkCamera() Camera {
	return NewCamera(
		256,
		144,
		100,
		10,
		2*runtime.NumCPU(),
		Point3{13, 2, 3},
		Point3{0, 0, 0},
		Vec3{0, 1, 0},
		20.0,
		0.1,
		10.0)
}

func BenchmarkRender(b *testing.B) {
	var (
		world      = randomScene()
		cam        = newBenchmarkCamera()
		mockWriter = writer{}
	)

	b.Run("render", func(b *testing.B) {
		cam.Render(world, mockWriter)
	})
}
