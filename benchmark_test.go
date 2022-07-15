package main

import "testing"

type writer struct{}

func (w writer) Write(p []byte) (int, error) {
	return len(p), nil
}

func BenchmarkProcess(b *testing.B) {
	const (
		w       = 400
		h       = int(w / ratio)
		samples = 10
	)
	var (
		world      = randomScene()
		cam        = newCamera()
		mockWriter = writer{}
	)

	b.Run("process", func(b *testing.B) {
		process(cam, world, w, h, samples, mockWriter)
	})
}
